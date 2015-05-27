package main

import (
    "github.com/spf13/viper";
    "flag";
    "fmt";
    "log";
    "os"
)

const VERSION = "0.0.1"

func main() {
    viper.SetConfigName("config")
    viper.SetDefault("BinDir", "/usr/local/bin")
    viper.AddConfigPath("$HOME/.dkenv")
    viper.ReadInConfig()

    version := flag.String("version", "", "Docker Version")
    list := flag.Bool("list", false, "Docker Version")
    apiVersion := flag.String("apiVersion", "", "API Version")

    var ver string

    flag.Parse()

     
    if *list {
        fmt.Println("Versions downloaded:")
        ListDownloadedVersions()
        os.Exit(0);
    }

    if (len(*apiVersion) > 0 || len(*version) > 0) {
        fmt.Println("version has value ", *version)
        fmt.Println("apiVersion has value ", *apiVersion)
        if len(*apiVersion) > 0 {
            var err error
            ver, err = ApiToVersion(*apiVersion)

            if err != nil {
                log.Fatal( err )
            }

            fmt.Println("For apiVersion ", *apiVersion, " using version ", ver)
        } else {
            ver = string(*version)
        }


        if (VersionDownloaded(ver)) {
            
        } else {
            GetDocker(ver, viper.GetString("BinDir"))
        }   
        SwitchVersion(ver, viper.GetString("BinDir"))
        
    } else {
        flag.Usage()
    }
}
