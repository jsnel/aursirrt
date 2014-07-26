package StorageCore

import (
	"github.com/joernweissenborn/AurSir4Go"
	PropertyGraph "github.com/joernweissenborn/propertygraph2go"
	uuid "github.com/nu7hatch/gouuid"
	"log"
)

type StorageCore struct {
	graph *PropertyGraph.PropertyGraph
	root  *PropertyGraph.Vertex
}

const (
	export_edge     = "EXPORT"
	import_edge     = "IMPORT"
	tag_edge        = "HAS_TAG"
	implements_edge = "IMPLEMENTS"
	awaiting_job_edge = "AWAITING_JOB"
	doing_job_edge = "DOING_JOB"
	listen_edge = "LISTEN"
)

func (sc *StorageCore) init() {

	sc.graph = PropertyGraph.New()

	sc.root = sc.graph.CreateVertex(generateUuid(), nil)

}

func (sc StorageCore) addImport(request AddImportRequest) (string, bool) {
	a := sc.graph.GetVertex(request.Id)
	if a == nil {
		log.Println("StorageCore Linking App as importer failed, app does exist:", request.Id)
		return "", false
	}
	kv := sc.registerKey(request.AppKey)
	vtags := make([]*PropertyGraph.Vertex, len(request.Tags))
	for i, t := range request.Tags {
		vtags[i] = sc.registerTag(t, kv)
	}

	return sc.registerImport(a, kv, vtags)
}

//registerImport creates or gets en import vertex and links it to the app, key and tags vertices
func (sc StorageCore) registerImport(app, key *PropertyGraph.Vertex, tags []*PropertyGraph.Vertex) (string, bool) {

	iv := sc.getImport(app, key)
	log.Println("StorageCore Linking App as importer")

	if iv == nil {
		iv = sc.graph.CreateVertex(generateUuid(), nil)
		sc.graph.CreateEdge(generateUuid(), import_edge, key, iv, nil)
		sc.graph.CreateEdge(generateUuid(), import_edge, iv, app, nil)
	} else {
		for _, te := range iv.Outgoing {
			if te.Label == tag_edge {
				sc.graph.RemoveEdge(te.Id)
			}
		}
	}

	for _, tag := range tags {
		sc.graph.CreateEdge(generateUuid(), tag_edge, tag, iv, nil)
	}

	f, _ := sc.isExported(iv)

	return iv.Id, f
}

//getImport gets an import object linked with a given app and key vertex
func (sc StorageCore) getImport(app, key *PropertyGraph.Vertex) *PropertyGraph.Vertex {

	//Imports are linked with keys via IMPORT edges from import to key vertex
	for _, ee := range key.Incoming {
		if ee.Label == import_edge {
			//Imports are linked with apps via IMPORT edges from app to import vertex
			for _, ae := range ee.Tail.Incoming {
				if ae.Label == import_edge && ae.Tail.Id == app.Id {
					return ae.Head
				}
			}
		}
	}

	return nil
}

func (sc StorageCore) addExport(expReq AddExportRequest) ExportAdded {
	//get the app vertex
	a := sc.graph.GetVertex(expReq.Id)
	if a == nil {
		log.Println("StorageCore Linking App as exporter failed, app does exist:", expReq.Id)
		return ExportAdded{}
	}

	//get or create the appkey vertex
	kv := sc.registerKey(expReq.AppKey)

	//prepare a slice for all tag vertices and then create or get them
	vtags := make([]*PropertyGraph.Vertex, len(expReq.Tags))
	for i, t := range expReq.Tags {
		vtags[i] = sc.registerTag(t, kv)
	}

	//register app key, returning the export_id
	return sc.registerExport(a, kv, vtags)
}

//registerExport creates or gets en export vertex and links it to the app, key and tags vertices. If the export already
//exists, all tags edges are deleted and rebuild
func (sc StorageCore) registerExport(app, key *PropertyGraph.Vertex, tags []*PropertyGraph.Vertex) ExportAdded {

	ev := sc.getExport(app, key)
	log.Println("StorageCore Linking App as exporter")
	connKeys := map[string]string{}

	if ev == nil {
		ev = sc.graph.CreateVertex(generateUuid(), nil)
		sc.graph.CreateEdge(generateUuid(), export_edge, key, ev, nil)
		sc.graph.CreateEdge(generateUuid(), export_edge, ev, app, nil)
	} else {
		for _, te := range ev.Outgoing {
			if te.Label == tag_edge {
				sc.graph.RemoveEdge(te.Id)
			}
		}
	}

	for _, tag := range tags {
		sc.graph.CreateEdge(generateUuid(), tag_edge, tag, ev, nil)
	}

	for _, imp := range key.Incoming {
		if imp.Label == import_edge {
			sc.linkExporterToListen(imp.Tail)
			if ok, _ := sc.isExported(imp.Tail); ok {
				connKeys[imp.Tail.Id] = sc.getImportApp(imp.Tail).Id
			}
		}
	}

	pending := []AurSir4Go.AurSirRequest{}
	for _, r := range sc.getPendingRequests(key){
		pending = append(pending,r.Properties.(AurSir4Go.AurSirRequest))
	}
	return ExportAdded{ev.Id, connKeys,pending}
}

func (sc StorageCore) getImportApp(iv *PropertyGraph.Vertex) *PropertyGraph.Vertex {
	for _, e := range iv.Incoming {
		if e.Label == import_edge {
			return e.Tail
		}
	}
	return nil
}

func (sc StorageCore) getImportKey(iv *PropertyGraph.Vertex) *PropertyGraph.Vertex {
	for _, e := range iv.Outgoing {
		if e.Label == import_edge {
			return e.Tail
		}
	}
	return nil
}

func (sc StorageCore) getTags(ImExport *PropertyGraph.Vertex) []*PropertyGraph.Vertex{
	tags := []*PropertyGraph.Vertex{}

	for _, te := range ImExport.Outgoing {
		if te.Label == tag_edge {
			tags = append(tags, te.Head)
		}
	}
	return tags
}
func (sc StorageCore) getTagNames(ImExport *PropertyGraph.Vertex) []string{
	tagnames := []string{}

	for _, t := range sc.getTags(ImExport) {
		tag, _ := t.Properties.(string)
		tagnames = append(tagnames, tag)
	}
	return tagnames
}

func (sc StorageCore) isCompatible(imp, exp *PropertyGraph.Vertex) bool{
	for _,tag := range sc.getTagNames(imp) {
		if !sc.hasTag(tag,exp){
			return false
		}
	}
	return true
}

func (sc StorageCore) hasTag(tag string, imporexp *PropertyGraph.Vertex) bool {
	for _, ele := range sc.getTagNames(imporexp) {
		if ele == tag {
			return true
		}
	}
	return false
}

func (sc StorageCore) getExporter(imp *PropertyGraph.Vertex) []*PropertyGraph.Vertex{
	//get the key
	var key *PropertyGraph.Vertex
		for _, ie := range imp.Outgoing {
		if ie.Label == import_edge {
			key = ie.Head
			break
		}
	}
	if key==nil {
		return nil
	}
	exporter := []*PropertyGraph.Vertex{}
	for _, e := range key.Incoming {
		if e.Label == export_edge {
			if sc.isCompatible(imp, e.Tail){
			exporter = append(exporter,e.Tail)
			}
		}
		}
	return exporter

}

func (sc StorageCore) isExported(imp *PropertyGraph.Vertex) (bool, *PropertyGraph.Vertex) {

	for _,e := range imp.Incoming {
		if e.Label == implements_edge {
			return true, e.Tail
		}
	}

	exps := sc.getExporter(imp)

	if len(exps)==0 {return false, nil}
	sc.graph.CreateEdge(generateUuid(),implements_edge,imp,exps[0],nil)
	return true, exps[0]
}
//getExportApp returns the app vertex to a gven export vertex
func (sc StorageCore) getExportApp(exp *PropertyGraph.Vertex) *PropertyGraph.Vertex {
	for _,e := range exp.Incoming {
		if e.Label == export_edge {
			return e.Tail
		}
	}
	return nil
}

//getExport gets an export object linked with a given app and key vertex
func (sc StorageCore) getExport(app, key *PropertyGraph.Vertex) *PropertyGraph.Vertex {

	//Exports are linked with keys via EXPORT edges from export to key vertex
	for _, ee := range key.Incoming {
		if ee.Label == export_edge {
			//Exports are linked with apps via EXPORT edges from app to export vertex
			for _, ae := range ee.Tail.Incoming {
				if ae.Label == export_edge && ae.Tail.Id == app.Id {
					return ae.Head
				}
			}
		}
	}

	return nil
}

func (sc StorageCore) registerTag(t string, k *PropertyGraph.Vertex) *PropertyGraph.Vertex {

	tv := sc.getTag(t, k)
	if tv == nil {
		tv = sc.graph.CreateVertex(generateUuid(), t)
		sc.graph.CreateEdge(generateUuid(), tag_edge, tv, k, nil)
	}
	return tv
}

func (sc StorageCore) getTag(t string, key *PropertyGraph.Vertex) *PropertyGraph.Vertex {

	for _, kv := range key.Outgoing {
		if kv.Label == tag_edge {
			tag, _ := kv.Head.Properties.(string)
			if tag == t {
				return kv.Head
			}
		}

	}

	return nil

}

func (sc StorageCore) getAllExports(app *PropertyGraph.Vertex) []*PropertyGraph.Vertex {
	a := []*PropertyGraph.Vertex{}
	for _, ee := range app.Outgoing {
		if ee.Label == export_edge {
			a = append(a, ee.Head)
		}
	}
	return a
}

func (sc StorageCore) getAllImports(app *PropertyGraph.Vertex) []*PropertyGraph.Vertex {
	a := []*PropertyGraph.Vertex{}
	for _, ie := range app.Outgoing {
		if ie.Label == import_edge {
			a = append(a, ie.Head)
		}
	}
	return a
}

func (sc StorageCore) getAllImExports(app *PropertyGraph.Vertex) []*PropertyGraph.Vertex {
	a := sc.getAllImports(app)
	for _, e := range sc.getAllExports(app) {
		a = append(a, e)
	}
	return a
}

func (sc StorageCore) registerApp(req RegisterAppRequest) {

	log.Println("Assinging ID to App:", req.AppName)

	sc.graph.CreateVertex(string(req.Id), req.AppName)

}

func (sc StorageCore) removeApp(req RemoveAppRequest) AppRemoved {
	log.Println("StorageCore Removing app", req.Id)
	discApps := map[string]string{}
	a := sc.graph.GetVertex(req.Id)
	for _, e := range sc.getAllImExports(a) {
		for _, impl := range e.Outgoing {
			if impl.Label == implements_edge {
				log.Println("StorageCore disconnecting app", req.Id)
				discApps[impl.Head.Id] = sc.getImportApp(impl.Head).Id
			}
		}
		sc.graph.RemoveVertex(e.Id)
	}

	sc.graph.RemoveVertex(req.Id)
	log.Println("StorageCore removal complete")
	return AppRemoved{discApps}
}

func (sc StorageCore) registerKey(k AurSir4Go.AppKey) *PropertyGraph.Vertex {

	log.Println("StorageCore Registering AppKey:", k.ApplicationKeyName)

	kv, f := sc.getKeyVertex(k.ApplicationKeyName)

	if f {
		log.Println("StorageCore Aborting register, key already known")
		return kv
	}

	kv = sc.graph.CreateVertex(generateUuid(), k)

	sc.graph.CreateEdge(generateUuid(), "KNOWN_APPKEY", kv, sc.root, nil)

	log.Println("StorageCore Key registered")

	return kv
}

func (sc StorageCore) getKeyVertex(k string) (*PropertyGraph.Vertex, bool) {
	for _, kv := range sc.root.Outgoing {
		key, _ := kv.Head.Properties.(AurSir4Go.AppKey)
		if key.ApplicationKeyName == k {
			return kv.Head, true
		}
	}
	return nil, false
}

func (sc StorageCore) addResult(arr AddResRequest) []string {
	req := arr.Req
	app := sc.graph.GetVertex(arr.AppId)
	if app == nil {
		log.Println("StorageCore error registering result, unknown app")
		return nil}

	key, f := sc.getKeyVertex(req.AppKeyName)
	if !f {log.Println("StorageCore error registering result, unknown key")
		return nil}

	exp:= sc.getExport(app,key)
	if exp == nil {log.Println("StorageCore error registering result, key not imported by app")
		return nil}

	job := sc.graph.GetVertex(arr.Req.Uuid)
	if job == nil {log.Println("StorageCore error registering result, job not found")
		return nil}


	if req.CallType == AurSir4Go.ONE2ONE || req.CallType == AurSir4Go.ONE2MANY {
		jobapp:=sc.getRequestApp(job)
		if jobapp != nil {
			sc.graph.RemoveVertex(job.Id)
			return []string{jobapp.Id}
		}
	}

	sc.graph.RemoveVertex(job.Id)

	if req.CallType == AurSir4Go.MANY2ONE || req.CallType == AurSir4Go.MANY2MANY {
		importer := []string{}

		for _,imp := range sc.getListener(exp){
			importer = append(importer,imp.Id)
		}
		return importer
	}

	return nil
}
// getRequestApp retrieves the app waiting for a job. Returns nil if job is MANY2..
func (sc StorageCore) getRequestApp(j *PropertyGraph.Vertex)*PropertyGraph.Vertex{
	for _,e := range j.Incoming{
		if e.Label == awaiting_job_edge {
			return sc.getImportApp(e.Tail)
		}
	}
	return nil
}


func (sc StorageCore) getPendingRequests(key *PropertyGraph.Vertex) []*PropertyGraph.Vertex{
	requests := []*PropertyGraph.Vertex{}
	for _,r := range sc.getRequests(key){
		if !sc.requestsProcessed(r) {
			requests = append(requests,r)
		}
	}
	return requests
}

func (sc StorageCore) requestsProcessed(r *PropertyGraph.Vertex) bool{
	for _, e := range r.Incoming{
		if e.Label == doing_job_edge {
			return true
		}
	}
	return false
}


func (sc StorageCore) getRequests(key *PropertyGraph.Vertex) []*PropertyGraph.Vertex{

	requests := []*PropertyGraph.Vertex{}

	for _,imp := range sc.getKeyImport(key){
		for _, e := range imp.Outgoing{
			if e.Label == awaiting_job_edge {
				requests = append(requests,e.Head)
			}
		}
	}
	return requests
}


func (sc StorageCore) getKeyImport(key *PropertyGraph.Vertex) []*PropertyGraph.Vertex{
	imports:= []*PropertyGraph.Vertex{}
	for _,e := range key.Incoming{
		if e.Label == import_edge {
			imports = append(imports, e.Tail)
		}
	}
	return imports
}

func (sc StorageCore) addFuncListen(req ListenRequest){
	imp := sc.graph.GetVertex(req.ImportId)
	if imp == nil {log.Println("StorageCore error registering listen,imported not found")
		return}



	for _,lf := range sc.getListeningFunctions(imp) {
		if lf.Properties.(string) == req.FuncName {
			return
		}
	}
	lf := sc.graph.CreateVertex(generateUuid(),req.FuncName)
	sc.graph.CreateEdge(generateUuid(),listen_edge,lf,imp,nil)
	for _,exp := range sc.getExporter(imp) {
		sc.graph.CreateEdge(generateUuid(),listen_edge,exp,lf,nil)
	}

}

func (sc StorageCore) getListeningFunctions(imp *PropertyGraph.Vertex) []*PropertyGraph.Vertex{
	lfs:= []*PropertyGraph.Vertex{}
	for _,e := range imp.Outgoing{
		if e.Label == listen_edge {
			lfs = append(lfs, e.Head)
		}
	}
	return lfs
}

func (sc StorageCore) linkExporterToListen(imp *PropertyGraph.Vertex){
	lfs := sc.getListeningFunctions(imp)
	for _,lf := range lfs {
		for _,e := range lf.Outgoing {
			sc.graph.RemoveEdge(e.Id)
		}
		for _,exp := range sc.getExporter(imp)		{
			sc.graph.CreateEdge(generateUuid(),listen_edge,exp,lf,nil)
		}

	}
}

func (sc StorageCore) getListener(exp *PropertyGraph.Vertex)[]*PropertyGraph.Vertex{
	listener := []*PropertyGraph.Vertex{}
	for _,e := range exp.Incoming{

		if e.Label == listen_edge {
			listener = append(listener, sc.getListeningApp(e.Tail))
		}
	}

	return listener
}

func (sc StorageCore) getListeningApp(lf *PropertyGraph.Vertex) *PropertyGraph.Vertex{
	for _,e := range lf.Incoming{
		if e.Label == listen_edge {
			return sc.getImportApp(e.Tail)
		}
	}
	return nil
}


func (sc StorageCore) addRequest(arr AddReqRequest) []string {
	req := arr.Req
	app := sc.graph.GetVertex(arr.AppId)
	if app == nil {
		log.Println("StorageCore error registering request, unknown app")
		return nil}

	key, f := sc.getKeyVertex(req.AppKeyName)
	if !f {log.Println("StorageCore error registering request, unknown key")
		return nil}

	imp:= sc.getImport(app,key)
	if imp == nil {log.Println("StorageCore error registering request, key not imported by app")
		return nil}

	rv := sc.graph.CreateVertex(req.Uuid,req)

	sc.graph.CreateEdge(generateUuid(),awaiting_job_edge,rv,imp,nil)

	if req.CallType == AurSir4Go.ONE2ONE || req.CallType == AurSir4Go.MANY2ONE {
		f ,export := sc.isExported(imp)
		if f {
			sc.graph.CreateEdge(generateUuid(),doing_job_edge,rv,export,nil)
			return []string{sc.getExportApp(export).Id}
		}
	}
	return nil
}

func generateUuid() string {
	Uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal("Failed to generate UUID")
		return ""
	}
	return Uuid.String()
}
