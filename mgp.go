package main

import (
  "errors"
  "os"
  "path/filepath"
)

// Process path and enqueues if ok for match checking
func (env *Env) processEntry(e *Entry) error {

  if e.ShouldSkip() {
    return filepath.SkipDir
  }

  if !e.ShouldProcess() {
    return nil
  }

  // hangs if the buffer is full
  env.sChan <- true
  // adds one goroutine to the wait group
  env.wg.Add(1)
  go func() {
    if env.matchContext {
      match, err := e.MatchAll()
      
      if err == nil && match != nil {
        env.msg.AddMatches(e.GetPath(), match)
      }
    } else {
      singleMatch, err := e.MatchFirst()

      if err == nil && singleMatch != nil {
        env.msg.AddMatch(e.GetPath(), singleMatch)
      }
    }

    // frees one position in the buffer
    <-env.sChan
    // signals goroutine finished
    env.wg.Done()
  }()

  return nil
}

func (env *Env) Run() {

  if _, err := os.Stat(env.startpath); os.IsNotExist(err) {
    env.msg.AddPathError(env.startpath, errors.New("path does not exists"))
    os.Exit(1)
  }

  // Traversing filepath
  filepath.Walk(env.startpath,

  func(pathname string, info os.FileInfo, err error) error {

    if *env.stopWalk {
      // If the termination is requested, the path Walking
      // stops and the function returns with an error
      return errors.New("user requested termination")
    }
    e := NewEntry(info, pathname, env)

    // Processes path in search of matches with the given
    // pattern or the excluded directories
    if err == nil {
      return env.processEntry(e)
    }

    // Checking permission and access errors
    env.msg.AddPathError(e.GetPath(), err)

    if e.node.IsDir() {
      return filepath.SkipDir
    }

    return nil 
  })

// Waits for goroutines to finish
  env.wg.Wait()
}
