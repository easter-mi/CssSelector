package cssSelector

func QueryFromCssExpression(node *html.Node,cssExpression string)(res []*html.Node){

} 

func preDeal(cssExpression string)string{
	fp:=regexp.MustCompile(`\s+`)
	p1:=fp.ReplaceAllLiteralString(pattern," ")

	fp1:=regexp.MustCompile(`\s?([+>,])\s?`)
	p1=fp1.ReplaceAllString(p1,"$1")
	p1=strings.Trim(p1," ")
	return p1
}
func getGroups(cssStdExpression string)[]string{
	return strings.Split(cssStdExpression,",")
}
func getAtomsAndJoiners(cssGroupExpression string)(atoms,atomJioners []string){
	// joinMap:=map[string]string{"+":"下一个",">":"子元素"," ":"内的"}
	atomeReg:=regexp.MustCompile(`([+> ])`)
	atoms:=atomeReg.Split(itemSelector,-1)
	atomJioners:=atomeReg.FindAllString(itemSelector,-1)
	return 
}

func queryFromCssAtom(node *html.Node,cssAtom string)(res []*html.Node){
		
}

