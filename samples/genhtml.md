func (g *Generator) GenerateIndex(t *template.Template, p []Post) error {

    // Set the export path and file name for the index.html generated code
    filePath := filepath.Join("index.html")
    f, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("error creating file %s: %v", filePath, err)
    }

    // Create a buffer to store the generated site temporarily
    w := bufio.NewWriter(f)
    // Use the Execute() function provided by the template/html package to generate the webpage
    if err := t.Execute(w, p); err != nil {
        return fmt.Errorf("error executing template %s : %v", filePath, err)
    }

    // Flush the generated file on the disk
    if err := w.Flush(); err != nil {
        return fmt.Errorf("error writing file %s: %v", filePath, err)
    }

    f.Close()
    return nil
}
