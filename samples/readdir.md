func (g *Generator) LoadPosts() ([]Post, error) {

    // Set the path for the root directory, where posts are located
    path := filepath.Join("posts")
    var postFolders []string

    // Open the root directory
    dir, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("error accessing directory %s: %v", path, err)
    }
    defer dir.Close()

    // Read all post directories. Each post is encapsulated in a different folder.
    files, err := dir.Readdir(-1)
    if err != nil {
        return nil, fmt.Errorf("error reading contents of directory %s: %v", path, err)
    }

    // Append each directory to an array
    for _, file := range files {
        if file.IsDir() && file.Name()[0] != '.' {
            postFolders = append(postFolders, filepath.Join(path, file.Name()))
        }
    }

    // Read all posts from the directories and parse them into an array of Post structs
    var posts []Post
    for _, folder := range postFolders {
        // Create a new Post struct
        post, err := newPost(folder)
        if err != nil {
            return nil, fmt.Errorf("error reading post contents %s: %v", folder, err)
        }
        posts = append(posts, post)
    }

    // Sort posts according to their ID
    sort.Slice(posts, func(i, j int) bool {
        return posts[i].Metadata.Id > posts[j].Metadata.Id
    })

    return posts, nil
}
