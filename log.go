package wbml

import "log"

const (
    normalText  = "\033[0m"
    boldText    = "\033[1m"
    blackText   = "\033[30;1m"
    redText     = "\033[31;1m"
    greenText   = "\033[32;1m"
    yellowText  = "\033[33;1m"
    blueText    = "\033[34;1m"
    magentaText = "\033[35;1m"
    cyanText    = "\033[36;1m"
    whiteText   = "\033[37;1m"
)
var (
    infoTag    = greenText + "INFO " + normalText
    warningTag = yellowText + "WARNING " + normalText
    debugTag   = blueText + "DEBUG " + normalText
    fatalTag   = redText + "FATAL " + normalText
)

func info(message string) {log.Print(infoTag + message)}
func warning(message string) {log.Print(warningTag + message)}
func debug(message string) {log.Print(debugTag + message)}
func fatal(message string) {log.Fatal(fatalTag + message)}

func infof(message string, values ...interface{}) {
    log.Printf(infoTag + message, values...)
}
func warningf(message string, values ...interface{}) {
    log.Printf(warningTag + message, values...)
}
func debugf(message string, values ...interface{}) {
    log.Printf(debugTag + message, values...)
}
func fatalf(message string, values ...interface{}) {
    log.Fatalf(fatalTag + message, values...)
}
