package mdson

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/drgo/core/ui"
)

//algorithm
// call evalBlock for each entry in root.children
// find all variable calls of the form {some key.some attrib};
// escape dot if it occurs in block or attrib name
//
// resolve using attrib map

func Eval(root BlockNode, options *ParserOptions) (BlockNode, error){
	ev := &evaluator{
		root: root ,
		UI: ui.NewUI(options.Debug),
	}
	if err:= ev.evalAllAttribs(root); err!=nil {
		return nil, err 
	}
	return ev.evalBlock(root)

}

type evaluator struct {
	root BlockNode
	ui.UI
}

var reCurelyBraces = regexp.MustCompile(`{([^{}]*)}`)

func (ev *evaluator) getAttribValue(att string) *string{
	a, ok := ev.root.Attribs()[att]
	if ok {
		return &a 
	}
	return nil 	
}


func (ev *evaluator) evalAttribRefs(s string) string{
	repl:= func (s string) string {
		if len(s) <3 {
			return fmt.Sprintf("<error: too short attribute '%s'>", s) 
		}
		s= strings.TrimSpace(s[1:len(s)-1]) //strip { and }
		if s== "" {
			return fmt.Sprintf("<error: too short attribute '%s'>", s) 
		}
		if expanded:=ev.getAttribValue(s); expanded != nil {
			return *expanded 
		}
		return fmt.Sprintf("<error: no such attribute '%s'>", s) 
	}	
	return reCurelyBraces.ReplaceAllStringFunc(s, repl)
}

func (ev *evaluator) evalLeaf(n Node)(Node, error) {
	//TODO: guard against evaluating empty, error what else?
	s:= ev.evalAttribRefs(n.Value())
	n.SetValue(s)	
	return n, nil
}

func (ev *evaluator) evalBlock(n BlockNode)(BlockNode, error) {
	ev.Log("evalBlock():", n)
	count:= len(n.Children())
	for i := 0; i < count; i++ {
		switch c:= n.NthChild(i).(type){
		case BlockNode:
			if _,err :=ev.evalBlock(c); err!=nil {
				return nil, err
			}
		//TODO: guard against evaluating errors etc
		default:
			if _, err :=ev.evalLeaf(n); err!=nil {
				return nil, err 
			}	
		}		
	}
	return n, nil  
}

func (ev *evaluator) evalAllAttribs(n BlockNode) error {
	for k,v := range n.Attribs(){
		s:= ev.evalAttribRefs(v)
		if s!=v {
			n.Attribs()[k]=s 
		}	
	}
	for _, c := range n.Children() {
		if c, ok:= c.(BlockNode); ok {
			if err:=ev.evalAllAttribs(c); err!= nil {
				return err 
			}
		}
	}	
	return nil  
}
