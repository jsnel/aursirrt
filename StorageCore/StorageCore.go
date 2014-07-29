package storagecore

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
	export_edge       = "EXPORT"
	import_edge       = "IMPORT"
	tag_edge          = "HAS_TAG"
	implements_edge   = "IMPLEMENTS"
	awaiting_job_edge = "AWAITING_JOB"
	doing_job_edge    = "DOING_JOB"
	listen_edge       = "LISTEN"
	callchain_edge = "CHAINCALL"
)

func (sc *StorageCore) init() {

	sc.graph = PropertyGraph.New()

	sc.root = sc.graph.CreateVertex(generateUuid(), nil)

}

func (sc StorageCore) addImport(request AddImportRequest) (string, bool) {
	app := sc.graph.GetVertex(request.Id)
	if app == nil {
		log.Println("StorageCore Linking App as importer failed, app does exist:", request.Id)
		return "", false
	}
	key := sc.registerKey(request.AppKey)
	if key == nil {
		log.Println("StorageCore Linking App as importer failed, key cannot be registered:", request.Id)
		return "", false
	}

	imp := sc.createImport(app, key)
	log.Println("StorageCore Linking App as importer with id",imp.Id)


	tags:=sc.registerTags(request.Tags,key)
	sc.linkTags(imp,tags)
	f, _ := sc.isExported(imp)

	return imp.Id, f
}
//updateExports removes all tags edges from the export vertex specified by the ExportId in the request and rebuilds
//them with the new tag set
func (sc StorageCore) updateImport(uir UpdateImportRequest) ImportAdded{
	log.Println("StorageCore updating import",uir.Req.ImportId)
	//grab the export
	imp := sc.graph.GetVertex(uir.Req.ImportId)
	if imp == nil {
		log.Println("StorageCore aborting export update, export does not exist")
		return ImportAdded{}
	}

	//stripe old tags
	for _, te := range imp.Outgoing {
		if te.Label == tag_edge {
			sc.graph.RemoveEdge(te.Id)
		}
	}
	//grab the key
	key := sc.getImportKey(imp)
	//Create or get new Tags:
	tags := sc.registerTags(uir.Req.Tags,key)
	sc.linkTags(imp,tags)

	sc.removeImplementor(imp)

	exported, _ := sc.isExported(imp)
	return ImportAdded{imp.Id, exported}
}
func (sc StorageCore) removeImplementor(imp *PropertyGraph.Vertex){
	for _, e := range imp.Incoming {
		if e.Label == implements_edge {
			sc.graph.RemoveEdge(e.Id)
		}
	}
}
//createImport creates an import object and links it with a given app and key vertex
func (sc StorageCore) createImport(app, key *PropertyGraph.Vertex) *PropertyGraph.Vertex {

	iv := sc.graph.CreateVertex(generateUuid(), nil)
	sc.graph.CreateEdge(generateUuid(), import_edge, key, iv, nil)
	sc.graph.CreateEdge(generateUuid(), import_edge, iv, app, nil)

	return iv
}

func (sc StorageCore) addExport(expReq AddExportRequest) ExportAdded {
	//get the app vertex
	app := sc.graph.GetVertex(expReq.Id)
	if app == nil {
		log.Println("StorageCore Linking App as exporter failed, app does exist:", expReq.Id)
		return ExportAdded{}
	}

	//get or create the appkey vertex
	key := sc.registerKey(expReq.AppKey)

	tags :=sc.registerTags(expReq.Tags,key)

	exp := sc.createExport(app, key)
	sc.linkTags(exp,tags)

	log.Println("StorageCore Linking App as exporter with export id",exp.Id)

	for _, imp := range sc.getKeyImport(key) {
			sc.linkExporterToListen(imp)
	}
	//no exports are going offline due to export adding
	return ExportAdded{exp.Id, sc.getExportedImportsForKey(key),map[string]string{}, sc.getPendingRequests(key)}
}


func (sc StorageCore) registerTags(Tags []string, key *PropertyGraph.Vertex)[]*PropertyGraph.Vertex{
	tags := make([]*PropertyGraph.Vertex, len(Tags))
	for i, t := range Tags {
		tags[i] = sc.registerTag(t, key)
	}
	return tags
}

//getExportedImportsForKey returns a map from ImportId to AppId, containing export app pairs that are currently exported
func (sc StorageCore) getExportedImportsForKey(key *PropertyGraph.Vertex) map[string]string{
	exported := map[string]string{}

	for _, imp := range key.Incoming {
			if cc,_:= sc.isChainCall(imp.Tail);!cc && imp.Label == import_edge  {
			if ok, _ := sc.isExported(imp.Tail); ok {
				exported[imp.Tail.Id] = sc.getImportApp(imp.Tail).Id
			}
		}
	}
	return exported
}

//getUnexportedImportsForKey returns a map from ImportId to AppId, containing export app pairs that are currently not
// exported
func (sc StorageCore) getUnexportedImportsForKey(key *PropertyGraph.Vertex) map[string]string{
	notexported := map[string]string{}

	for _, imp := range key.Incoming {
		if imp.Label == import_edge {
			if ok, _ := sc.isExported(imp.Tail); !ok {
				notexported[imp.Tail.Id] = sc.getImportApp(imp.Tail).Id
			}
		}
	}
	return notexported
}

func (sc StorageCore) linkTags(imExport *PropertyGraph.Vertex,tags []*PropertyGraph.Vertex){
	for _, tag := range tags {
		sc.graph.CreateEdge(generateUuid(), tag_edge, tag, imExport, nil)
	}
}

//updateExports removes all tags edges from the export vertex specified by the ExportId in the request and rebuilds
//them with the new tag set
func (sc StorageCore) updateExport(uer UpdateExportRequest) ExportAdded {
	log.Println("StorageCore updating export",uer.Req.ExportId)
	//grab the export
	exp := sc.graph.GetVertex(uer.Req.ExportId)
	if exp == nil {
		log.Println("StorageCore aborting export update, export does not exist")
		return ExportAdded{}
	}

	//stripe old tags
	for _, te := range exp.Outgoing {
		if te.Label == tag_edge {
			sc.graph.RemoveEdge(te.Id)
		}
	}
	//grab the key
	key := sc.getExportKey(exp)

	//Create or get new Tags:
	tags := sc.registerTags(uer.Req.Tags,key)
	sc.linkTags(exp,tags)
	sc.removeALlKeyImportImplementor(key)
	connected:= sc.getExportedImportsForKey(key)
	return ExportAdded{uer.Req.ExportId, connected, sc.getUnexportedImportsForKey(key), sc.getPendingRequests(key)}
}

func (sc StorageCore) removeALlKeyImportImplementor(key *PropertyGraph.Vertex){
	for _, imp :=range sc.getKeyImport(key) {
		sc.removeImplementor(imp)
	}
}

func (sc StorageCore) getImportApp(iv *PropertyGraph.Vertex) *PropertyGraph.Vertex {
	for _, e := range iv.Incoming {
		if e.Label == import_edge {
			return e.Tail
		}
	}
	return nil
}

//getImportKey delivers app key vertex associated with a given import vertex
func (sc StorageCore) getImportKey(importVertex *PropertyGraph.Vertex) *PropertyGraph.Vertex {
	for _, e := range importVertex.Outgoing {
		if e.Label == import_edge {
			return e.Head
		}
	}
	return nil
}

//getExportKey delivers appkey vertex associated with a given export vertex
func (sc StorageCore) getExportKey(exportVertex *PropertyGraph.Vertex) *PropertyGraph.Vertex {
	for _, e := range exportVertex.Outgoing {
		if e.Label == export_edge {
			return e.Head
		}
	}
	return nil
}

func (sc StorageCore) getTags(ImExport *PropertyGraph.Vertex) []*PropertyGraph.Vertex {
	tags := []*PropertyGraph.Vertex{}

	for _, te := range ImExport.Outgoing {
		if te.Label == tag_edge {
			tags = append(tags, te.Head)
		}
	}
	return tags
}
func (sc StorageCore) getTagNames(ImExport *PropertyGraph.Vertex) []string {
	tagnames := []string{}

	for _, t := range sc.getTags(ImExport) {
		tag, _ := t.Properties.(string)
		tagnames = append(tagnames, tag)
	}
	return tagnames
}

func (sc StorageCore) isCompatible(imp, exp *PropertyGraph.Vertex) bool {
	for _, tag := range sc.getTagNames(imp) {
		if !sc.hasTag(tag, exp) {
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

func (sc StorageCore) getExporter(imp *PropertyGraph.Vertex) []*PropertyGraph.Vertex {
	//get the key
	var key *PropertyGraph.Vertex
	for _, ie := range imp.Outgoing {
		if ie.Label == import_edge{
			key = ie.Head
			break
		}
	}
	if key == nil {
		return nil
	}
	exporter := []*PropertyGraph.Vertex{}
	for _, e := range key.Incoming {
		if e.Label == export_edge {
			if sc.isCompatible(imp, e.Tail) {
				exporter = append(exporter, e.Tail)
			}
		}
	}
	return exporter

}

func (sc StorageCore) isExported(imp *PropertyGraph.Vertex) (bool, *PropertyGraph.Vertex) {

	for _, e := range imp.Incoming {
		if e.Label == implements_edge {
			return true, e.Tail
		}
	}

	exps := sc.getExporter(imp)

	if len(exps) == 0 {
		return false, nil
	}
	sc.graph.CreateEdge(generateUuid(), implements_edge, imp, exps[0], nil)
	return true, exps[0]
}

//getExportApp returns the app vertex to a gven export vertex
func (sc StorageCore) getExportApp(exp *PropertyGraph.Vertex) *PropertyGraph.Vertex {
	for _, e := range exp.Incoming {
		if e.Label == export_edge {
			return e.Tail
		}
	}
	return nil
}

//createExport creates an export object linked with a given app and key vertex
func (sc StorageCore) createExport(app, key *PropertyGraph.Vertex) *PropertyGraph.Vertex {

	//Exports are linked with keys via EXPORT edges from export to key vertex
	ev := sc.graph.CreateVertex(generateUuid(), nil)
	sc.graph.CreateEdge(generateUuid(), export_edge, key, ev, nil)
	sc.graph.CreateEdge(generateUuid(), export_edge, ev, app, nil)

	return ev
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
		if cc,_:= sc.isChainCall(ie.Head);ie.Label == import_edge && !cc {
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
			log.Println("StorageCore disconnecting app", impl.Head.Id)

			if iscc,_:= sc.isChainCall(impl.Head);impl.Label == implements_edge && !iscc {
				log.Print("StorageCore disconnecting app", req.Id)
				log.Println(" from", impl.Head.Id)
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

	//log.Println("StorageCore AppKey Hash =", k.Hash())

	kv := sc.getKeyVertex(k.ApplicationKeyName)

	if kv != nil {
		log.Println("StorageCore Aborting register, key already known")
		return kv
	}

	kv = sc.graph.CreateVertex(generateUuid(), k)

	sc.graph.CreateEdge(generateUuid(), "KNOWN_APPKEY", kv, sc.root, nil)

	log.Println("StorageCore Key registered")

	return kv
}

func (sc StorageCore) getKeyVertex(k string) (*PropertyGraph.Vertex) {
	for _, kv := range sc.root.Outgoing {
		key, _ := kv.Head.Properties.(AurSir4Go.AppKey)
		if key.ApplicationKeyName == k {
			return kv.Head
		}
	}
	return nil
}

func (sc StorageCore) addResult(arr AddResRequest) ResRegistered {
	req := arr.Req
	app := sc.graph.GetVertex(arr.AppId)
	if app == nil {
		log.Println("StorageCore error registering result, unknown app")
		return ResRegistered{}
	}


	exp := sc.graph.GetVertex(req.ExportId)
	if exp == nil {
		log.Println("StorageCore error registering result, key not imported by app")
		return ResRegistered{}
	}

	job := sc.graph.GetVertex(arr.Req.Uuid)
	if job == nil {
		log.Println("StorageCore error registering result, job not found")
		return ResRegistered{}
	}

	ischaincall , ccv:= sc.hasChainCall(job)
	var chaincall AurSir4Go.ChainCall
	chaincallimportid := ""
	if ischaincall {
		chaincall, _ = ccv.Properties.(AurSir4Go.ChainCall)
		for _,e := range ccv.Incoming {
			if e.Label == awaiting_job_edge {
				chaincallimportid= e.Tail.Id
			}
		}

	}

	if req.CallType == AurSir4Go.ONE2ONE || req.CallType == AurSir4Go.ONE2MANY {
		jobapp := sc.getRequestApp(job)
		log.Println("StorageCore error registering result, requesting app not found")
		if jobapp != nil {
			sc.graph.RemoveVertex(job.Id)
			return ResRegistered{[]string{jobapp.Id},ischaincall,chaincall,chaincallimportid}
		}
	}

	sc.graph.RemoveVertex(job.Id)

	if req.CallType == AurSir4Go.MANY2ONE || req.CallType == AurSir4Go.MANY2MANY {
		importer := []string{}

		for _, imp := range sc.getListener(exp) {
			importer = append(importer, imp.Id)
		}
		return ResRegistered{importer,ischaincall,chaincall,chaincallimportid}
	}

	return  ResRegistered{[]string{},ischaincall,chaincall,chaincallimportid}
}

// getRequestApp retrieves the app waiting for a job. Returns nil if job is MANY2..
func (sc StorageCore) getRequestApp(j *PropertyGraph.Vertex) *PropertyGraph.Vertex {
	for _, e := range j.Incoming {
		if e.Label == awaiting_job_edge {
			return sc.getImportApp(e.Tail)
		}
	}
	return nil
}

func (sc StorageCore) getPendingRequests(key *PropertyGraph.Vertex) []AurSir4Go.AurSirRequest {

	requests := []AurSir4Go.AurSirRequest{}
	for _, r := range sc.getRequests(key) {
		if !sc.requestsProcessed(r) {
			requests = append(requests, r.Properties.(AurSir4Go.AurSirRequest))
		}
	}
	return requests
}

func (sc StorageCore) requestsProcessed(r *PropertyGraph.Vertex) bool {
	for _, e := range r.Incoming {
		if e.Label == doing_job_edge {
			return true
		}
	}
	return false
}

func (sc StorageCore) getRequests(key *PropertyGraph.Vertex) []*PropertyGraph.Vertex {

	requests := []*PropertyGraph.Vertex{}

	for _, imp := range sc.getKeyImport(key) {
		for _, e := range imp.Outgoing {
			if e.Label == awaiting_job_edge {
				requests = append(requests, e.Head)
			}
		}
	}
	return requests
}

func (sc StorageCore) getKeyImport(key *PropertyGraph.Vertex) []*PropertyGraph.Vertex {
	imports := []*PropertyGraph.Vertex{}
	for _, e := range key.Incoming {
		if cc,_ :=sc.isChainCall(e.Tail) ;e.Label == import_edge && !cc {
			imports = append(imports, e.Tail)
		}
	}
	return imports
}

func (sc StorageCore) addFuncListen(req ListenRequest) {
	imp := sc.graph.GetVertex(req.ImportId)
	if imp == nil {
		log.Println("StorageCore error registering listen,imported not found")
		return
	}

	for _, lf := range sc.getListeningFunctions(imp) {
		if lf.Properties.(string) == req.FuncName {
			return
		}
	}
	lf := sc.graph.CreateVertex(generateUuid(), req.FuncName)
	sc.graph.CreateEdge(generateUuid(), listen_edge, lf, imp, nil)
	for _, exp := range sc.getExporter(imp) {
		sc.graph.CreateEdge(generateUuid(), listen_edge, exp, lf, nil)
	}

}

func (sc StorageCore) getListeningFunctions(imp *PropertyGraph.Vertex) []*PropertyGraph.Vertex {
	lfs := []*PropertyGraph.Vertex{}
	for _, e := range imp.Outgoing {
		if e.Label == listen_edge {
			lfs = append(lfs, e.Head)
		}
	}
	return lfs
}

func (sc StorageCore) linkExporterToListen(imp *PropertyGraph.Vertex) {
	lfs := sc.getListeningFunctions(imp)
	for _, lf := range lfs {
		for _, e := range lf.Outgoing {
			sc.graph.RemoveEdge(e.Id)
		}
		for _, exp := range sc.getExporter(imp) {
			sc.graph.CreateEdge(generateUuid(), listen_edge, exp, lf, nil)
		}

	}
}

func (sc StorageCore) getListener(exp *PropertyGraph.Vertex) []*PropertyGraph.Vertex {
	listener := []*PropertyGraph.Vertex{}
	for _, e := range exp.Incoming {

		if e.Label == listen_edge {
			listener = append(listener, sc.getListeningApp(e.Tail))
		}
	}

	return listener
}

func (sc StorageCore) getListeningApp(lf *PropertyGraph.Vertex) *PropertyGraph.Vertex {
	for _, e := range lf.Incoming {
		if e.Label == listen_edge {
			return sc.getImportApp(e.Tail)
		}
	}
	return nil
}

func (sc StorageCore) addRequest(arr AddReqRequest) []string {
	req := arr.Req

	imp := sc.graph.GetVertex(req.ImportId)
	if imp == nil {
		log.Println("StorageCore error registering request, key not imported by app")
		return nil
	}

	rv := sc.graph.GetVertex(req.Uuid)
	if rv == nil {
		rv = sc.graph.CreateVertex(req.Uuid, req)
		sc.graph.CreateEdge(generateUuid(), awaiting_job_edge, rv, imp, nil)
	}


	if req.CallType == AurSir4Go.ONE2ONE || req.CallType == AurSir4Go.MANY2ONE {
		f, export := sc.isExported(imp)
		if f {
			sc.graph.CreateEdge(generateUuid(), doing_job_edge, rv, export, nil)
			return []string{sc.getExportApp(export).Id}
		}
	}
	return nil
}



func (sc StorageCore) addCallChain(accr AddCallChainRequest)[]string{
	req := accr.Req
	app := sc.graph.GetVertex(accr.AppId)
	if app == nil {
		log.Println("StorageCore error registering request, unknown app")
		return nil
	}


	imp := sc.graph.GetVertex(req.OriginRequest.ImportId)
	if imp == nil {
		log.Println("StorageCore error registering request, key not imported by app")
		return nil
	}

	rv := sc.graph.CreateVertex(req.OriginRequest.Uuid, req.OriginRequest)

	sc.graph.CreateEdge(generateUuid(), awaiting_job_edge, rv, imp, nil)
	prev := rv
	for _, call := range req.CallChain {
		call.ChainCallId =generateUuid()
		log.Println("StorageCore creating ChainCall with Id:",call.ChainCallId)
		cv := sc.graph.CreateVertex(call.ChainCallId, call)
		tags := sc.registerTags(call.Tags,sc.getKeyVertex(call.AppKeyName))
		sc.linkTags(cv,tags)
		key := sc.getKeyVertex(call.AppKeyName)
		sc.graph.CreateEdge(generateUuid(), import_edge, cv, app, nil)
		sc.graph.CreateEdge(generateUuid(), import_edge, key, cv, nil)
		sc.graph.CreateEdge(generateUuid(), awaiting_job_edge, cv, cv, nil)
		sc.graph.CreateEdge(generateUuid(), callchain_edge, cv, prev, nil)
		prev = cv
	}

	if req.FinalImportId != ""{
		fimp := sc.graph.GetVertex(req.FinalImportId)
		if fimp == nil {
			log.Println("StorageCore error registering final request, unknown import")
			return nil
		}

		frv := sc.graph.CreateVertex(req.FinalCall.ChainCallId,req.FinalCall)
		sc.graph.CreateEdge(generateUuid(), awaiting_job_edge, frv, fimp, nil)
		sc.graph.CreateEdge(generateUuid(), callchain_edge, frv, prev, nil)

	}

	if req.OriginRequest.CallType == AurSir4Go.ONE2ONE || req.OriginRequest.CallType == AurSir4Go.MANY2ONE {
		f, export := sc.isExported(imp)
		if f {
			sc.graph.CreateEdge(generateUuid(), doing_job_edge, rv, export, nil)
			return []string{sc.getExportApp(export).Id}
		}
	}
	return nil

}

func (sc StorageCore) isChainCall(req *PropertyGraph.Vertex) (bool, *PropertyGraph.Vertex) {
	for _,e := range req.Incoming {
		if e.Label == awaiting_job_edge && e.Tail == e.Head {
			return true, e.Head
		}
	}
	return false, nil
}

func (sc StorageCore) hasChainCall(req *PropertyGraph.Vertex) (bool, *PropertyGraph.Vertex) {
	for _,e := range req.Outgoing {
		if e.Label == callchain_edge {
			return true, e.Head
		}
	}
	return false, nil
}

func generateUuid() string {
	Uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal("Failed to generate UUID")
		return ""
	}
	return Uuid.String()
}
