package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

/* Read the source Latex file */
func readLatexFile(latexFilePath string) string {
	file, err := os.Open(latexFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	
	// replace line-breakers with blankspace in the content
	//content := strings.Replace(string(data), "\n", " ", -1)
	return string(data)
}


// TODO: deal with the micros defined by "\newcommand{}{}"
func processMacros(latexContent string) string{
	return latexContent
}

/* 
1. get the contents between "\begin{document}" and "\end{document}"
2. insert <p> and </p> tags at proper places 
3. remove "\title{}" "\date{}" "\maketitle" "\bibliographystyle{}" and "\bibliography{}" if there is any 
4. replace macros defined by "\newcommand{}{}"
*/
func getBody(latexContent string) string {

	// handle the macros
	latexContent = processMacros(latexContent)

	// "(?s)" is the flag to turn on the DOTALL mode, which allows "." to match '\n' also.
	reg := regexp.MustCompile(`(?s)(\\begin{document})(.*)(\\end{document})`)
	body := reg.FindStringSubmatch(latexContent)[2]
	
	body = regexp.MustCompile(`(?s)\\title{.*?}` +
									`|\\date{.*?}` +
									`|\\maketitle` +
									`|\\bibliographystyle{.*?}` +
									`|\\bibliography{.*?}`).ReplaceAllString(body, "")

	// replace (more-than) double line-breakers with "<p>" tag	
	newBody := ""
	for _,para := range regexp.MustCompile(`\n\n+`).Split(body, -1) {
		if len(para) != 0 {
			newBody = newBody + "<p>" + para + "</p>"
		}
	}
	// Each remaining return is converted to " " (i.e. a blankspace)
    newBody = regexp.MustCompile(`\n`).ReplaceAllString(newBody, " ")

	return newBody
}

// handle "itemize" and "enumerate" environment that take no paramters
func processEnumItem(body string) string {
	enumItemTags := map[string]string{"itemize": "ul",
									"enumerate": "ol",
									"item": "li"} 

	reg := regexp.MustCompile(`\\begin{itemize}` +
								`|\\end{itemize}` +
								`|\\begin{enumerate}` + 
								`|\\end{enumerate}` +
								`|\\item`)							
	textBlocks := reg.Split(body,-1)
	controlTags := reg.FindAllString(body, -1)

	newBody := textBlocks[0]
	for i,tag := range controlTags{

		// handle "\item" tage seperately
		if tag == "\\item" {
			tagType := "item"
			newBody += "<" + enumItemTags[tagType] + ">" + textBlocks[i+1] + "</" + enumItemTags[tagType] + ">"
			continue
		}

		parsed := regexp.MustCompile(`\\|{|}`).Split(tag,-1)
		tagType := parsed[1]
		tagValue :=  parsed[2]
		if tagType == "begin"{
			newBody += "<" + enumItemTags[tagValue] + ">" 
		}else{
			newBody += "</" + enumItemTags[tagValue] + ">" 
		}
		 
		newBody += textBlocks[i+1]
		
	}
	return newBody
}

func processTheorems(body string) string {
	theoremTags := map[string]string{"theorem": "blockquote",
									"title": "b"} 

	var counter int32 = 1

	labelTable := make(map[string]string)


	reg := regexp.MustCompile(`\\begin{theorem}` +
								`|\\end{theorem}` +
								`|\\label{theorem:.*?}` +
								`|\\begin{lemma}` + 
								`|\\end{lemma}` +
								`|\\label{lemma:.*?}`)							
	textBlocks := reg.Split(body,-1)
	controlTags := reg.FindAllString(body, -1)

	newBody := textBlocks[0]
	for i,tag := range controlTags{
		// Skip labels
		if strings.Contains(tag, "\\label{"){ 
			newBody += textBlocks[i+1]
			continue 
		}

		htmlTheorem := ""
		label := ""
		// if it is the openning
		if strings.Contains(tag, "\\begin{"){ 

			parsed := regexp.MustCompile(`{|}`).Split(tag,-1)
			theoremType := parsed[1]
			theoremIndex := fmt.Sprint(counter)

			// if this theorem/lemma is bond with a label
			if i+1 < len(controlTags) && strings.Contains(controlTags[i+1], "\\label{"){
				labelParser := regexp.MustCompile(`\\|{|}`).Split(controlTags[i+1],-1)
				label =  labelParser[2]
				// add an record to the labelTable
				labelTable[label] = theoremIndex
			}
			// increase the counter by 1
			counter++

			htmlTheorem = "<" + theoremTags["theorem"] +  " id=\"" + label + "\">"
			htmlTheorem += "<" + theoremTags["title"] + ">" + strings.Title(theoremType) + " " + theoremIndex  +  "</" + theoremTags["title"] + ">" 
		}else{
			// if it is the closing
			htmlTheorem = "</" + theoremTags["theorem"] + ">"
		}

		newBody += htmlTheorem
		newBody += textBlocks[i+1]
	}	

	// replace the \ref{} calls to theorems/lemmata properly
	// This block of code is identical to the one used in processSections(). Maybe make it a seperate function later when optimizing the code 
	reg = regexp.MustCompile(`\\ref{.*?}`)					
	textBlocks = reg.Split(newBody,-1)
	controlTags = reg.FindAllString(newBody, -1)
	newBody = textBlocks[0]
	for i,tag := range controlTags{
		parsed := regexp.MustCompile(`{|}`).Split(tag,-1)
		refLabel := parsed[1]
		htmlRef := ""
		if  index, exist := labelTable[refLabel]; exist {
			htmlRef = "<a href=\"#" + refLabel + "\">" + index + "</a>"
		}else{
			// keep the original "\ref{.*?}" without touch
			htmlRef = tag
		}
		newBody = newBody + htmlRef
		newBody += textBlocks[i+1]
	}

	// HTML does not allow ":" appear in the "id" for anchor, so replace it by "_"
	newBody = regexp.MustCompile(`theorem:`).ReplaceAllString(newBody, "theorem_")
	newBody = regexp.MustCompile(`lemma:`).ReplaceAllString(newBody, "lemma_")
	return newBody
}


func processSections(body string) string {
	sectionTags := map[string]string{"section": "h2", "subsection": "h3"} 

	var sectionCounter int32 = 1
	var subsectionCounter int32 = 1

	labelTable := make(map[string]string)

	reg := regexp.MustCompile(`\\section{.*?}` +
								`|\\subsection{.*?}` + 
								`|\\label{section:.*?}`)
								
	textBlocks := reg.Split(body,-1)
	controlTags := reg.FindAllString(body, -1)

	newBody := textBlocks[0]
	for i,tag := range controlTags{
		// Skip labels
		if strings.Contains(tag, "\\label{"){ 
			newBody += textBlocks[i+1]
			continue 
		}

		parsed := regexp.MustCompile(`\\|{|}`).Split(tag,-1)
		tagType := parsed[1]
		tagValue :=  parsed[2]

		sectionIndex := ""
		
		if tagType == "section" {
			sectionIndex = fmt.Sprint(sectionCounter)
			sectionCounter++
			subsectionCounter = 1
		}else if tagType == "subsection"{

			sectionIndex = fmt.Sprint(sectionCounter) + "." + fmt.Sprint(subsectionCounter)
			subsectionCounter++
		}

		htmlSection := ""
		// if this section/subsection is bond with a label
		if i+1 < len(controlTags) && strings.Contains(controlTags[i+1], "\\label{"){
			labelParser := regexp.MustCompile(`\\|{|}`).Split(controlTags[i+1],-1)
			label :=  labelParser[2]

			// add an record to the labelTable
			labelTable[label] = sectionIndex

			htmlSection = "<" + sectionTags[tagType] + " id=\"" + label + "\">" + sectionIndex + " " +  tagValue + "</" + sectionTags[tagType] + ">"
		}else{
			htmlSection = "<" + sectionTags[tagType] + ">" + sectionIndex + " " +  tagValue + "</" + sectionTags[tagType] + ">"
		} 

		newBody = newBody + "\n" +  htmlSection + "\n"
		newBody += textBlocks[i+1]
	}
	
	// replace the \ref{} calls to sections/subsections properly
	reg = regexp.MustCompile(`\\ref{.*?}`)					
	textBlocks = reg.Split(newBody,-1)
	controlTags = reg.FindAllString(newBody, -1)
	newBody = textBlocks[0]
	for i,tag := range controlTags{
		parsed := regexp.MustCompile(`{|}`).Split(tag,-1)
		refLabel := parsed[1]
		htmlRef := ""
		if  index, exist := labelTable[refLabel]; exist {
			htmlRef = "<a href=\"#" + refLabel + "\">" + index + "</a>"
		}else{
			// keep the original "\ref{.*?}" without touch
			htmlRef = tag
		}
		newBody = newBody + htmlRef
		newBody += textBlocks[i+1]
	}


	// HTML does not allow ":" appear in the "id" for anchor, so replace it by "_"
	newBody = regexp.MustCompile(`section:`).ReplaceAllString(newBody, "section_")
	return newBody
}


func writeToFile(content string, filePath string){
	file, err := os.Create(filePath)
    if err != nil {
        fmt.Println(err)
        return
	}
	defer file.Close()

    file.WriteString(content)
} 

func main() {
	content := readLatexFile("Resource/test.tex")
	
	body := getBody(content) 

	//fmt.Println(body)
	//fmt.Println(string(data))
	body = processSections(body)

	body = processTheorems(body)

	body = processEnumItem(body)

	
	// break lines properly to easy the reading
	body = regexp.MustCompile(`</p><p>`).ReplaceAllString(body, "</p>\n\n<p>")
	body = regexp.MustCompile(`<ul>`).ReplaceAllString(body, "\n<ul>")
	body = regexp.MustCompile(`<ol>`).ReplaceAllString(body, "\n<ol>")
	body = regexp.MustCompile(`</ol>`).ReplaceAllString(body, "\n</ol>\n")
	body = regexp.MustCompile(`</ul>`).ReplaceAllString(body, "\n</ul>\n")
	body = regexp.MustCompile(`<li>`).ReplaceAllString(body, "\n<li>")


    writeToFile("[latexpage]\n" + body, "test.txt")
	
	//fmt.Printf("%s", body)

}
