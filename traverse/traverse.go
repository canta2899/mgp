package traverse

import (
  "errors"
  "os"
  "path/filepath"
)

// Process path and enqueues if ok for match checking
func (en *Env) ProcessEntry(e *Entry) error {

  if e.ShouldSkip() {
    return filepath.SkipDir
  }

  if !e.ShouldProcess() {
    return nil
  }

  // hangs if the buffer is full
  en.Schan <- true
  // adds one goroutine to the wait group
  en.Wg.Add(1)
  go func() {
    if en.MatchContext {
      match, err := e.MatchAll()
      
      if err == nil && match != nil {
        en.Msg.AddMatches(e.GetPath(), match)
      }
    } else {
      singleMatch, err := e.MatchFirst()

      if err == nil && singleMatch != nil {
        en.Msg.AddMatch(e.GetPath(), singleMatch)
      }
    }

    // frees one position in the buffer
    <-en.Schan
    // signals goroutine finished
    en.Wg.Done()
  }()

  return nil
}

func (en *Env) Run() {

  if _, err := os.Stat(en.StartPath); os.IsNotExist(err) {
    en.Msg.AddPathError(en.StartPath, errors.New("path does not exists"))
    os.Exit(1)
  }

  // Traversing filepath
  filepath.Walk(en.StartPath,

  func(pathname string, info os.FileInfo, err error) error {

    if *en.StopWalk {
      // If the termination is requested, the path Walking
      // stops and the function returns with an error
      return errors.New("user requested termination")
    }
    e := NewEntry(info, pathname, en)

    // Processes path in search of matches with the given
    // pattern or the excluded directories
    if err == nil {
      return en.ProcessEntry(e)
    }

    // Checking permission and access errors
    en.Msg.AddPathError(e.GetPath(), err)

    if e.Node.IsDir() {
      return filepath.SkipDir
    }

    return nil 
  })

// Waits for goroutines to finish
  en.Wg.Wait()
}
