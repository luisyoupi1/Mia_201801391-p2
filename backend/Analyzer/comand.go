package analyzer_test

import (
	"P2/Functions"
	"P2/Utilities"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)


func Command(input string) {
	
	// Verificar si el input está vacío
	if input == "" {
		return // No hacer nada si el input está vacío
	}

	comando := input
	input = strings.ToLower(input)
	switch {
	case strings.HasPrefix(input, "mkdisk"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleMKDISKCommand(comando)

	case strings.HasPrefix(input, "rmdisk"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleRMDISKCommand(comando)

	case strings.HasPrefix(input, "fdisk"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleFDISKCommand(comando)

	case strings.HasPrefix(input, "mount"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleMOUNTCommand(comando)

	case strings.HasPrefix(input, "unmount"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleUNMOUNTCommand(comando)

	case strings.HasPrefix(input, "mkfs"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleMKFSCommand(comando)

	case strings.HasPrefix(input, "login"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleLOGINCommand(comando)

	case strings.HasPrefix(input, "logout"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleLOGOUTCommand()

	case strings.HasPrefix(input, "mkgrp"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleMKGRPCommand(comando)

	case strings.HasPrefix(input, "rmgrp"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleRMGRPCommand(comando)

	case strings.HasPrefix(input, "mkusr"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleMKUSRCommand(comando)

	case strings.HasPrefix(input, "rmusr"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleRMUSRCommand(comando)

	case strings.HasPrefix(input, "mkfile"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleMKFILECommand(comando)

	case strings.HasPrefix(input, "cat"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleCATCommand(comando)

	case strings.HasPrefix(input, "remove"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleREMOVECommand(comando)

	case strings.HasPrefix(input, "edit"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleEDITCommand(comando)

	case strings.HasPrefix(input, "rename"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleRENAMECommand(comando)

	case strings.HasPrefix(input, "mkdir"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleMKDIRCommand(comando)

	case strings.HasPrefix(input, "copy"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleCOPYCommand(comando)

	case strings.HasPrefix(input, "move"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleMOVECommand(comando)

	case strings.HasPrefix(input, "find"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleFINDCommand(comando)

	case strings.HasPrefix(input, "chown"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleCHOWNCommand(comando)

	case strings.HasPrefix(input, "chgrp"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleCHGRPCommand(comando)

	case strings.HasPrefix(input, "chmod"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleCHMODCommand(comando)

	case strings.HasPrefix(input, "pause"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handlePAUSECommand()

	case strings.HasPrefix(input, "execute"):
		handleEXECUTECommand(comando)

	case strings.HasPrefix(input, "rep"):
		AddText(">>>>>>>>>>>>>>>>>>>>" + comando+"")
		handleREPCommand(comando)

	case strings.HasPrefix(input, "#"):
		//Ignora las sentencias del lado derecho
	default:
		AddText("Comando no reconocido:"+ comando+"")
	}
}

var (
	size        = flag.Int("size", 0, "Tamaño")
	fit         = flag.String("fit", "", "Ajuste")
	unit        = flag.String("unit", "", "Unidad")
	type_       = flag.String("type", "", "Tipo")
	driveletter = flag.String("driveletter", "", "Busqueda")
	name        = flag.String("name", "", "Nombre")
	delete      = flag.String("delete", "", "Eliminar")
	add         = flag.Int("add", 0, "Añadir/Quitar")
	path        = flag.String("path", "", "Directorio")
	id          = flag.String("id", "", "ID")
	fs          = flag.String("fs", "", "FDISK")
	ruta        = flag.String("ruta", "", "Ruta")
	user        = flag.String("user", "", "Usuario")
	pass        = flag.String("pass", "", "Password")
	grp         = flag.String("grp", "", "Group")
	r           = flag.Bool("r", false, "Rewrite")
	cont        = flag.String("cont", "", "Cont")
	destino     = flag.String("destino", "", "Destino")
	ugo         = flag.String("ugo", "", "UGO")
	file        = flag.String("file", "", "File to process")
	flagN       = flag.Bool("error", false, "Flag not found")
)

func AddText(text string)  {
	utilities_test.Resultados.WriteString(text+"\n")
}

/* -------------------------------------------------------------------------- */
/*                           APLICACION DE COMANDOS                           */
/* -------------------------------------------------------------------------- */
func handleMKDISKCommand(input string) {

	flag.Parse()
	functions_test.ProcessMKDISK(input, size, fit, unit, flagN)

	if *flagN {
		*flagN = false
		return
	}

	// validate size > 0
	if *size <= 0 {
		AddText("Error: Size must be greater than 0")
		return
	}

	// validate fit equals to b/w/f
	if *fit != "b" && *fit != "f" && *fit != "w" {
		AddText("Error: Fit must be (bf/ff/wf)")
		return
	}

	// validate unit equals to k/m
	if *unit != "k" && *unit != "m" {
		AddText("Error: Unit must be (k/m)")
		return
	}
	
	AddText("--------------------------------------------------------------------------")
	AddText("                       MKDISK: PARAMETROS CORRECTOS                       ")
	AddText("--------------------------------------------------------------------------")

	// Create the file
	functions_test.CreateBinFile(size, fit, unit)
	*size = 0
	*fit = ""
	*unit = ""
}

func handleRMDISKCommand(input string) {
	flag.Parse()
	functions_test.ProcessRMDISK(input, driveletter, flagN)

	if *flagN {
		*flagN = false
		return
	}

	// validate driveletter be a letter and not empty
	if !functions_test.ValidDriveLetter(*driveletter) {
		AddText("Error: DriveLetter debe ser una letra")
		return
	} else if len(*driveletter) == 0 {
		AddText("Error: DriveLetter es un campo obligatorio")
		return
	}

	AddText("--------------------------------------------------------------------------")
	AddText("                       RMDISK: PARAMETROS CORRECTOS                       ")
	AddText("--------------------------------------------------------------------------")

	functions_test.DeleteBinFile(driveletter)
	*driveletter = ""
}

func handleFDISKCommand(input string) {
	flag.Parse()
	functions_test.ProcessFDISK(input, size, driveletter, name, unit, type_, fit, delete, add, path, flagN)

	if *flagN {
		*flagN = false
		return
	}

	//Obligatorio cuando no existe la particion
	// validate size > 0
	if *size <= 0 && *delete != "full" && *add == 0 {
		AddText("Error: Size must be greater than 0")
		return
	}

	// validate driveletter be a letter and not empty
	if !functions_test.ValidDriveLetter(*driveletter) {
		AddText("Error: DriveLetter must be a letter")
		return
	} else if len(*driveletter) == 0 {
		AddText("Error: DriveLetter cannot be empty")
		return
	}

	// validate fit equals to b/w/f
	if *fit != "b" && *fit != "f" && *fit != "w" {
		AddText("Error: Fit must be (BF/FF/WF)")
		return
	}

	// validate unit equals to b/k/m
	if *unit != "b" && *unit != "k" && *unit != "m" {
		AddText("Error: Unit must be (B/K/M)")
		return
	}

	//AddText("ADD", *add)
	// validate type equals to P/E/L
	if *type_ != "p" && *type_ != "e" && *type_ != "l" && *delete != "full" && *add == 0 {
		AddText("Error: Type must be (P/E/L)")
		return
	}

	if *delete != "" {
		if *delete != "full" {
			AddText("Error: Delete must be full")
			return
		}
		if *name == "" && *path == "" {
			AddText("Error: you need path and name to delete")
			return
		}
	}

	AddText("--------------------------------------------------------------------------")
	AddText("                        FDISK: PARAMETROS CORRECTOS                       ")
	AddText("--------------------------------------------------------------------------")

	functions_test.CRUD_Partitions(size, driveletter, name, unit, type_, fit, delete, add, path)
	*size = 0
	*driveletter = ""
	*name = ""
	*unit = ""
	*type_ = ""
	*fit = ""
	*delete = ""
	*add = 0
	*path = ""
}

func handleMOUNTCommand(input string) {
	flag.Parse()
	functions_test.ProcessMOUNT(input, driveletter, name, flagN)

	if *flagN {
		*flagN = false
		return
	}

	// validate driveletter be a letter and not empty
	if !functions_test.ValidDriveLetter(*driveletter) {
		AddText("Error: DriveLetter must be a letter")
		return
	} else if len(*driveletter) == 0 {
		AddText("Error: DriveLetter cannot be empty")
		return
	}

	AddText("--------------------------------------------------------------------------")
	AddText("                        MOUNT: PARAMETROS CORRECTOS                       ")
	AddText("--------------------------------------------------------------------------")

	functions_test.MountPartition(driveletter, name)
	*driveletter = ""
	*name = ""
}

func handleUNMOUNTCommand(input string) {
	flag.Parse()
	functions_test.ProcessUNMOUNT(input, id, flagN)

	if *flagN {
		*flagN = false
		return
	}

	if *id == "" {
		AddText("Error: Id es un campo obligatorio")
	}

	letra := string((*id)[0])
	AddText("DISCO:" + letra)

	if !functions_test.ValidDriveLetter(letra) {
		AddText("Error: ID")
		AddText("Error: DISCO INCORRECTO")
		return
	}

	numero := string((*id)[1])
	AddText("PARTICION:" + numero+"")

	if !utilities_test.EsNumero(numero) {
		AddText("Error: ID")
		AddText("Error: PARTICION INCORRECTA")
		return
	}

	AddText("CODIGO:" + string((*id)[2]) + string((*id)[3])+"")

	if string((*id)[2]) != "0" && string((*id)[3]) != "2" {
		AddText("Error: ID")
		AddText("Error: CODIGO INCORRECTO")
		return
	}

	AddText("--------------------------------------------------------------------------")
	AddText("                       UNMOUNT: PARAMETROS CORRECTOS                      ")
	AddText("--------------------------------------------------------------------------")

	functions_test.UNMOUNT_Partition(id)
	*id = ""
}

func handleMKFSCommand(input string) {
	flag.Parse()
	functions_test.ProcessMKFS(input, id, type_, fs, flagN)

	if *flagN {
		*flagN = false
		return
	}

	if *id == "" {
		AddText("Error: id es obligatorio")
	}

	if *fs != "2fs" && *fs != "3fs" {
		AddText("Error: fs debe ser 2fs o 3fs")
	}

	AddText("--------------------------------------------------------------------------")
	AddText("                         MKFS: PARAMETROS CORRECTOS                       ")
	AddText("--------------------------------------------------------------------------")

	functions_test.MKFS(id, type_, fs)
	*id = ""
	*type_ = ""
	*fs = ""
}

/* -------------------------------------------------------------------------- */
/*                         ADMINISTRACION DE USUARIOS                         */
/* -------------------------------------------------------------------------- */
func handleLOGINCommand(input string) {
	flag.Parse()
	functions_test.ProcessLOGIN(input, user, pass, id, flagN)

	if *user == "" || *pass == "" || *id == "" {
		AddText("--------------------------------------------------------------------------")
		AddText("                       LOGIN: PARAMETROS INCOMPLETOS                      ")
		AddText("--------------------------------------------------------------------------")

		return
	}

	if *flagN {
		*flagN = false
		return
	}

	AddText("--------------------------------------------------------------------------")
	AddText("                        LOGIN: PARAMETROS CORRECTOS                       ")
	AddText("--------------------------------------------------------------------------")

	functions_test.LOGIN(user, pass, id)

	*user = ""
	*pass = ""
	*id = ""
}

func handleLOGOUTCommand() {
	functions_test.LOGOUT()
}

func handleMKGRPCommand(input string) {
	flag.Parse()
	functions_test.ProcessMKGRP(input, name, flagN)

	if *name == "" {
		AddText("Error: el campo name no puede estar vacio")
		return
	}

	if *flagN {
		*flagN = false
		return
	}

	AddText("--------------------------------------------------------------------------")
	AddText("                        MKGRP: PARAMETROS CORRECTOS                       ")
	AddText("--------------------------------------------------------------------------")

	functions_test.MKGRP(name)
	*name = ""
}

func handleRMGRPCommand(input string) {
	flag.Parse()
	functions_test.ProcessRMGRP(input, name, flagN)

	if *name == "" {
		AddText("Error: el campo name no puede estar vacio")
		return
	}

	if *flagN {
		*flagN = false
		return
	}

	AddText("--------------------------------------------------------------------------")
	AddText("                        RMGRP: PARAMETROS CORRECTOS                       ")
	AddText("--------------------------------------------------------------------------")

	functions_test.RMGRP(name)
	*name = ""
}

func handleCHGRPCommand(input string) {
	flag.Parse()
	functions_test.ProcessCHGRP(input, user, grp, flagN)

	if *flagN {
		*flagN = false
		return
	}

	if *user == "" || *grp == "" {
		AddText("Error: campos incompletos")
		return
	}

	AddText("--------------------------------------------------------------------------")
	AddText("                        CHGRP: PARAMETROS CORRECTOS                       ")
	AddText("--------------------------------------------------------------------------")
	functions_test.CHGRP(user, grp)
	*user = ""
	*grp = ""
}

func handleMKUSRCommand(input string) {
	flag.Parse()
	functions_test.ProcessMKUSR(input, user, pass, grp, flagN)

	if len(*user) > 10 {
		AddText("Error: user no puede ser mayor a 10 caracteres")
		return
	}
	if len(*pass) > 10 {
		AddText("Error: password no puede ser mayor a 10 caracteres")
		return
	}
	if len(*grp) > 10 {
		AddText("Error: grupo no puede ser mayor a 10 caracteres")
		return
	}

	if *user == "" || *pass == "" || *grp == "" {
		AddText("Error: campos incompletos")
		return
	}

	if *flagN {
		*flagN = false
		return
	}

	AddText("--------------------------------------------------------------------------")
	AddText("                        MKUSR: PARAMETROS CORRECTOS                       ")
	AddText("--------------------------------------------------------------------------")

	functions_test.MKUSR(user, pass, grp)

	*user = ""
	*pass = ""
	*grp = ""
}

func handleRMUSRCommand(input string) {
	flag.Parse()
	functions_test.ProcessRMUSR(input, user, flagN)

	if *user == "" {
		AddText("Error: user no puede estar vacio")
		return
	}

	if *flagN {
		*flagN = false
		return
	}

	AddText("--------------------------------------------------------------------------")
	AddText("                        RMUSR: PARAMETROS CORRECTOS                       ")
	AddText("--------------------------------------------------------------------------")

	functions_test.RMUSR(user)

	*user = ""
}

/* -------------------------------------------------------------------------- */
/*                         ADMINISTRACION DE CARPETAS                         */
/* -------------------------------------------------------------------------- */
func handleMKDIRCommand(input string) {
	flag.Parse()
	functions_test.ProcessMKDIR(input, path, r, flagN)

	if *path == "" {
		AddText("Error: path no puede estar vacio")
		return
	}

	if *flagN {
		*flagN = false
		return
	}

	// AddText("Path: " + *path)
	// fmt.Print("r: ")
	// AddText(*r)

	AddText("--------------------------------------------------------------------------")
	AddText("                        MKDIR: PARAMETROS CORRECTOS                       ")
	AddText("--------------------------------------------------------------------------")

	functions_test.MKDIR(path, r)

	*path = ""
	*r = false
}

func handleMKFILECommand(input string) {
	flag.Parse()
	functions_test.ProcessMKFILE(input, path, r, size, cont, flagN)

	if *flagN {
		*flagN = false
		return
	}

	if *path == "" {
		AddText("Error: path no puede estar vacio")
		return
	}

	if *cont != "" {
		// Verificar si la ruta existe
		if _, err := os.Stat(*cont); err == nil {
			AddText("--------------------------------------------------------------------------")
			AddText("                        MKFILE: LA RUTA EXISTE                            ")
			AddText("--------------------------------------------------------------------------")
			AddText("La ruta existe en el sistema.")
		} else if os.IsNotExist(err) {
			AddText("Error: La ruta no existe en el sistema.")
			return
		} else {
			AddText("Error: No se logro verificar la ruta:"+ err.Error()+"")
			return
		}
	}

	if *size < 0 {
		AddText("Error: size negativo")
		*size = 0
		return
	}

	AddText("--------------------------------------------------------------------------")
	AddText("                       MKFILE: PARAMETROS CORRECTOS                       ")
	AddText("--------------------------------------------------------------------------")

	functions_test.MKFILE(path, r)
	*path = ""
	*r = false
	*size = 0
	*cont = ""
}

// A PARTIR DE AQUI FALTA EL FLAGN
func handleCATCommand(input string) {
	flag.Parse()
	functions_test.ProcessCAT(input, file)
}

func handleREMOVECommand(input string) {
	flag.Parse()
	functions_test.ProcessREMOVE(input, path, flagN)

	if *flagN {
		*flagN = false
		return
	}

	if *path == "" {
		AddText("Error: path vacio")
		return
	}

	functions_test.REMOVE(path)
	*path = ""
}

func handleEDITCommand(input string) {
	flag.Parse()
	functions_test.ProcessEDIT(input, path, cont)
}

func handleRENAMECommand(input string) {
	flag.Parse()
	functions_test.ProcessRENAME(input, path, name)
	if *flagN {
		*flagN = false
		return
	}

	if *path == "" {
		AddText("Error: path vacio")
		return
	}
	if *name == "" {
		AddText("Error: name vacio")
		return
	}

	functions_test.RENAME(path, name)
	*path = ""
	*name = ""
}

func handleCOPYCommand(input string) {
	flag.Parse()
	functions_test.ProcessCOPY(input, path, destino)

	if *path == "" || *destino == "" {
		AddText("Error: campos incompletos")
		return
	}

	functions_test.COPY(path, destino)
	*path = ""
	*destino = ""
}

func handleMOVECommand(input string) {
	flag.Parse()
	functions_test.ProcessMOVE(input, path, destino)
	if *path == "" || *destino == "" {
		AddText("Error: campos incompletos")
		return
	}

	functions_test.MOVE(path, destino)
	*path = ""
	*destino = ""
}

func handleFINDCommand(input string) {
	flag.Parse()
	functions_test.ProcessFIND(input, path, destino)
}

func handleCHOWNCommand(input string) {
	flag.Parse()
	functions_test.ProcessCHOWN(input, path, user, r)
}

func handleCHMODCommand(input string) {
	flag.Parse()
	functions_test.ProcessCHMOD(input, path, ugo, r)
}

/* -------------------------------------------------------------------------- */
/*                            COMANDOS AUXILIARES                             */
/* -------------------------------------------------------------------------- */
func handlePAUSECommand() {
	AddText("Presione cualquier tecla para continuar...")
	fmt.Scanln() // Espera a que el usuario presione Enter
	AddText("Continuando la ejecución...")
}

func handleEXECUTECommand(input string) {
	flag.Parse()
	functions_test.ProcessExecute(input, path, flagN)

	if *flagN {
		*flagN = false
		return
	}

	if *path == "" {
		AddText("Error: Path cannot be empty")
		return
	}
	// Open bin file
	file, err := utilities_test.OpenFile(*path)
	if err != nil {
		return
	}

	// Close bin file
	defer file.Close()

	// Crea un nuevo scanner para leer el archivo
	scanner := bufio.NewScanner(file)

	// Itera sobre cada línea del archivo
	for scanner.Scan() {
		linea := scanner.Text() // Lee la línea actual
		//AddText(linea)
		Command(linea)
	}

	// Verifica si hubo algún error durante la lectura
	if err := scanner.Err(); err != nil {
		AddText("Error al leer el archivo:"+ err.Error()+"")
	}
	*path = ""
}

/* -------------------------------------------------------------------------- */
/*                                  REPORTES                                  */
/* -------------------------------------------------------------------------- */
func handleREPCommand(input string) {
	flag.Parse()
	functions_test.ProcessREP(input, name, path, id, ruta, flagN)

	if *flagN {
		*flagN = true
		return
	}

	if *name == "" || *path == "" || *id == "" {
		AddText("Error: incomplete statements")
		return
	}

	letra := string((*id)[0])
	AddText("DISCO:" + letra+"")

	if !functions_test.ValidDriveLetter(letra) {
		AddText("Error: ID")
		AddText("Error: DISCO INCORRECTO")
		return
	}

	numero := string((*id)[1])
	AddText("PARTICION:" + numero+"")

	if !utilities_test.EsNumero(numero) {
		AddText("Error: ID")
		AddText("Error: PARTICION INCORRECTA")
		return
	}

	AddText("CODIGO:" + string((*id)[2]) + string((*id)[3])+"")

	if string((*id)[2]) != "0" && string((*id)[3]) != "2" {
		AddText("Error: ID")
		AddText("Error: CODIGO INCORRECTO")
		return
	}

	functions_test.GenerateReports(name, path, id, ruta)
}
