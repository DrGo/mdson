package mdson

import (
	"fmt"
	"regexp"
	"strings"
)

//algorithm
// call evalBlock for each entry in root.children
// find all variable calls of the form {some key.some attrib};
// escape dot if it occurs in block or attrib name
//
// resolve using attrib map

func (doc *Document) eval()  error{
	// assert(doc!=nil, "doc is nil in eval()")
	if err:= doc.evalAllAttribs(doc.root); err!=nil {
		return err 
	}
	doc.ctx.Log("eval()==> #sections", len(doc.root.Children()))
	n, err:= doc.evalBlock(doc.root)
	if err!=nil {
		return err 
	}
	doc.ctx.Log("eval()==> #sections of evaluted root", len(n.Children()))
	 doc.root = n 
	doc.ctx.Log("eval()==> #sections", len(doc.root.Children()))
	return nil
}

var reCurelyBraces = regexp.MustCompile(`{([^{}]*)}`)

func (doc *Document) getAttribValue(att string) *string{
	a, ok := doc.Attribs()[att]
	if ok {
		return &a 
	}
	return nil 	
}


func (doc *Document) evalAttribRefs(s string) string{
	repl:= func (s string) string {
		if len(s) <3 {
			return fmt.Sprintf("<error: too short attribute '%s'>", s) 
		}
		s= strings.TrimSpace(s[1:len(s)-1]) //strip { and }
		if s== "" {
			return fmt.Sprintf("<error: too short attribute '%s'>", s) 
		}
		if expanded:=doc.getAttribValue(s); expanded != nil {
			return *expanded 
		}
		return fmt.Sprintf("<error: no such attribute '%s'>", s) 
	}	
	return reCurelyBraces.ReplaceAllStringFunc(s, repl)
}

func (doc *Document) evalLeaf(n Node)(Node, error) {
	//TODO: guard against evaluating empty, error what else?
	s:= doc.evalAttribRefs(n.Value())
	doc.ctx.Log("**************** evalLeaf(): " + s)
	n.SetValue(s)	
	return n, nil
}

func (doc *Document) evalBlock(n BlockNode)(BlockNode, error) {
	doc.ctx.Log("evalBlock() start:", n.Key())
	count:= len(n.Children())
	for i := 0; i < count; i++ {
		switch c:= n.NthChild(i).(type){
		case BlockNode:
			_,err :=doc.evalBlock(c)  
			if err!=nil {
				return nil, err
			}
			// n.UpdateChild(i, en)
		//TODO: guard against evaluating errors etc
		default:
			if _, err :=doc.evalLeaf(c); err!=nil {
				return nil, err 
			}	
		}		
	}
	return n, nil  
}

func (doc *Document) evalAllAttribs(n BlockNode) error {
	for k,v := range doc.Attribs(){
		s:= doc.evalAttribRefs(v)
		if s!=v {
			doc.Attribs()[k]=s 
		}	
	}
	return nil  
}
