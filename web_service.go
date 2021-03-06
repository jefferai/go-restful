package restful

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"log"
)

// WebService holds a collection of Route values that bind a Http Method + URL Path to a function.
type WebService struct {
	rootPath       string
	pathExpr       *pathExpression // cached compilation of rootPath as RegExp
	routes         []Route
	produces       []string
	consumes       []string
	pathParameters []*Parameter
	filters        []FilterFunction
	documentation  string
}

// Path specifies the root URL template path of the WebService.
// All Routes will be relative to this path.
func (w *WebService) Path(root string) *WebService {
	w.rootPath = root
	compiled, err := newPathExpression(root)
	if err != nil {
		log.Fatalf("[restful] Invalid path:%s because:%v", root, err)
	}
	w.pathExpr = compiled
	return w
}

// Param adds a PathParameter to document parameters used in the root path.
func (w *WebService) Param(parameter *Parameter) *WebService {
	if w.pathParameters == nil {
		w.pathParameters = []*Parameter{}
	}
	w.pathParameters = append(w.pathParameters, parameter)
	return w
}

// PathParameter creates a new Parameter of kind Path for documentation purposes.
// It is initialized as required with string as its DataType.
func (w *WebService) PathParameter(name, description string) *Parameter {
	p := &Parameter{&ParameterData{Name: name, Description: description, Required: true, DataType: "string"}}
	p.bePath()
	return p
}

// QueryParameter creates a new Parameter of kind Query for documentation purposes.
// It is initialized as not required with string as its DataType.
func (w *WebService) QueryParameter(name, description string) *Parameter {
	p := &Parameter{&ParameterData{Name: name, Description: description, Required: false, DataType: "string"}}
	p.beQuery()
	return p
}

// BodyParameter creates a new Parameter of kind Body for documentation purposes.
// It is initialized as required without a DataType.
func (w *WebService) BodyParameter(name, description string) *Parameter {
	p := &Parameter{&ParameterData{Name: name, Description: description, Required: true}}
	p.beBody()
	return p
}

// HeaderParameter creates a new Parameter of kind (Http) Header for documentation purposes.
// It is initialized as not required with string as its DataType.
func (w *WebService) HeaderParameter(name, description string) *Parameter {
	p := &Parameter{&ParameterData{Name: name, Description: description, Required: false, DataType: "string"}}
	p.beHeader()
	return p
}

// FormParameter creates a new Parameter of kind Form (using application/x-www-form-urlencoded) for documentation purposes.
// It is initialized as required with string as its DataType.
func (w *WebService) FormParameter(name, description string) *Parameter {
	p := &Parameter{&ParameterData{Name: name, Description: description, Required: false, DataType: "string"}}
	p.beForm()
	return p
}

// Route creates a new Route using the RouteBuilder and add to the ordered list of Routes.
func (w *WebService) Route(builder *RouteBuilder) *WebService {
	builder.copyDefaults(w.produces, w.consumes)
	w.routes = append(w.routes, builder.Build())
	return w
}

// Method creates a new RouteBuilder and initialize its http method
func (w *WebService) Method(httpMethod string) *RouteBuilder {
	return new(RouteBuilder).servicePath(w.rootPath).Method(httpMethod)
}

// Produces specifies that this WebService can produce one or more MIME types.
// Http requests must have one of these values set for the Accept header.
func (w *WebService) Produces(contentTypes ...string) *WebService {
	w.produces = contentTypes
	return w
}

// Consumes specifies that this WebService can consume one or more MIME types.
// Http requests must have one of these values set for the Content-Type header.
func (w *WebService) Consumes(accepts ...string) *WebService {
	w.consumes = accepts
	return w
}

// Routes returns the Routes associated with this WebService
func (w WebService) Routes() []Route {
	return w.routes
}

// RootPath returns the RootPath associated with this WebService. Default "/"
func (w WebService) RootPath() string {
	return w.rootPath
}

// PathParameters return the path parameter names for (shared amoung its Routes)
func (w WebService) PathParameters() []*Parameter {
	return w.pathParameters
}

// Filter adds a filter function to the chain of filters applicable to all its Routes
func (w *WebService) Filter(filter FilterFunction) *WebService {
	w.filters = append(w.filters, filter)
	return w
}

// Doc is used to set the documentation of this service.
func (w *WebService) Doc(plainText string) {
	w.documentation = plainText
}

// Documentation returns it.
func (w WebService) Documentation() string {
	return w.documentation
}

/*
	Convenience methods
*/

// HEAD is a shortcut for .Method("HEAD").Path(subPath)
func (w *WebService) HEAD(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(w.rootPath).Method("HEAD").Path(subPath)
}

// GET is a shortcut for .Method("GET").Path(subPath)
func (w *WebService) GET(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(w.rootPath).Method("GET").Path(subPath)
}

// POST is a shortcut for .Method("POST").Path(subPath)
func (w *WebService) POST(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(w.rootPath).Method("POST").Path(subPath)
}

// PUT is a shortcut for .Method("PUT").Path(subPath)
func (w *WebService) PUT(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(w.rootPath).Method("PUT").Path(subPath)
}

// PATCH is a shortcut for .Method("PATCH").Path(subPath)
func (w *WebService) PATCH(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(w.rootPath).Method("PATCH").Path(subPath)
}

// DELETE is a shortcut for .Method("DELETE").Path(subPath)
func (w *WebService) DELETE(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(w.rootPath).Method("DELETE").Path(subPath)
}
