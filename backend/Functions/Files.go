package functions_test

import (
	"P2/Structs"
	"P2/Utilities"
	"encoding/binary"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	Padre         structs_test.Content
	padreBusqueda int
)

func AddText(text string)  {
	utilities_test.Resultados.WriteString(text + "\n")
}

//?               ADMINISTRACION DE CARPETAS, ARCHIVOS Y PERMISOS
/* -------------------------------------------------------------------------- */
/*                                COMANDO MKDIR                               */
/* -------------------------------------------------------------------------- */
func ProcessMKDIR(input string, path *string, r *bool, flagN *bool) {
	padreBusqueda = 0
	Padre.B_inodo = -1
	copy(Padre.B_name[:], "")
	flags := strings.Split(input, "-")
	for _, i := range flags {
		if i == "r" {
			*r = true
		}
		f := strings.Split(i, "=")
		if f[0] == "path" {
			*path = f[1]

			if strings.Contains(f[1], " ") {
				tieneComillas := strings.Split(*path, "\"")
				if len(tieneComillas)-1 == 0 {
					*path = `"` + f[1] + `"`
				}

			}
		}
	}

	AddText("--------------------------------------------------------------------------")
	AddText("                   		MKDIR: PROCESANDO	 	                       ")
	AddText(fmt.Sprint("\n" + *path + "\n\n"))
	AddText("--------------------------------------------------------------------------")

	re := regexp.MustCompile(`-(\w+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]

		switch flagName {
		case "r":
			*r = true
		case "path":
		default:
		}
	}
}

func MKDIR(path *string, r *bool) {
	/* -------------------------------------------------------------------------- */
	/*                  COMPROBAMOS SI HAY UNA SESSION EXISTENTE                  */
	/* -------------------------------------------------------------------------- */
	if !session {
		AddText("--------------------------------------------------------------------------")
		AddText("                   MKDIR: NO HAY UNA SESION INICIADA                      ")
		AddText("--------------------------------------------------------------------------")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepaths := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepaths)
	if err != nil {
		AddText("Error opening disk file:"+ err.Error())
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		AddText("Error reading MBR:"+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL INODO 0                            */
	/* -------------------------------------------------------------------------- */

	var Inode0 structs_test.Inode
	if err := utilities_test.ReadObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs_test.Inode{})))); err != nil {
		AddText("Error reading inode:"+ err.Error())
		return
	}


	/* -------------------------------------------------------------------------- */
	/*                           OBTENEMOS LA RUTA PADRE                          */
	/* -------------------------------------------------------------------------- */
	tieneComillas := strings.Split(*path, "\"")
	if len(tieneComillas)-1 != 0 {
		if len(tieneComillas)-1 == 1 {
			*path = tieneComillas[0]
		} else {
			*path = tieneComillas[1]
		}
	}
	rutaPadre := filepath.Dir(*path)
	AddText("Ruta original")
	AddText(*path)
	AddText("Ruta padre")
	AddText(rutaPadre)
	Carpetas := strings.Split(*path, "/")
	tieneArchivo := strings.Split(Carpetas[len(Carpetas)-1], ".")
	if (len(tieneArchivo) - 1) != 0 {
		AddText("Error: para crear archivos debes usar MKFILE")
		return
	}
	nuevaCarpeta := Carpetas[len(Carpetas)-1]
	partes := strings.Split(rutaPadre, "/")
	partes = partes[1:]
	//AddText("Elementos en ruta padre")
	//AddText(len(partes) - 1)
	carpetaCreada := false
	/* -------------------------------------------------------------------------- */
	/*                     RECORREMOS LOS BLOQUES DEL INODO 0                     */
	/* -------------------------------------------------------------------------- */
	//AddText("Bloques del inodo 0:")
	ultimo := 0
	root := false
	padreExiste := false
	for cont, i := range Inode0.I_block {
		if len(partes)-1 == 0 {
			//AddText("root es true")
			root = true
		}
		if i == -1 {
			ultimo = int(cont - 1)
			break
		}
		//AddText(i)
		if !root {
			existe := BuscarRuta(partes, i, 0)
			if existe {
				AddText("Existe la ruta padre")
				padreExiste = true
			}
		}
	}

	if root {
		existe := false
		for _, i := range Inode0.I_block {
			if i == -1 {
				break
			}
			// print("Buscando en el inodo 0 el bloque ")
			// AddText(i)
			existe = BuscarEspacioEnRoot(nuevaCarpeta, i)
			
			if existe {
				break
			}
		}
		if !existe {
			AddText("Creando nuevo inodo y bloque")
			// BlockCounter++
			CrrSuperblock.S_blocks_count++
			Inode0.I_block[ultimo+1] = CrrSuperblock.S_blocks_count
			CrearFolderBlock(file, CrrSuperblock.S_blocks_count, nuevaCarpeta)
			AddText("Actualizando inodo 0")
			if err := utilities_test.WriteObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs_test.Inode{})))); err != nil {
				AddText("Error reading inode:"+ err.Error())
				return
			}

		}
		carpetaCreada = true
	}

	if padreExiste && !carpetaCreada {
		AddText("Creando carpeta desde padre")
		CreandoCamino(Padre.B_inodo, nuevaCarpeta, file, partes)
		carpetaCreada = true
	}

	if *r && !carpetaCreada {
		if string(Padre.B_name[:]) != "" {
			AddText("creando a partir de carpetas existentes")
			AddText(fmt.Sprintf("Encontrado -> B_inode: %d B_name: %s\n", Padre.B_inodo, Padre.B_name))
			CreandoCamino(Padre.B_inodo, nuevaCarpeta, file, partes)
		} else {
			AddText("Creando todas las carpetas")
			CreandoCamino(0, nuevaCarpeta, file, partes)
		}
		carpetaCreada = true
	}
	if carpetaCreada {
		AddText("--------------------------------------------------------------------------")
		AddText(fmt.Sprintf("                MKDIR: CARPETA %s CREADA CORRECTAMENTE", nuevaCarpeta))
		AddText("--------------------------------------------------------------------------")
	} else {
		AddText("Error: No se logro crear la carpeta")
	}
	/* -------------------------------------------------------------------------- */
	/*                         ACTUALIZAMOS EL SUPERBLOQUE                        */
	/* -------------------------------------------------------------------------- */
	if err := utilities_test.WriteObject(file, &CrrSuperblock, int64(TempMBR.Mbr_particion[indexSB].Part_start)); err != nil {
		AddText("Error reading superblock:"+ err.Error())
		return
	}
}

/* -------------------------------------------------------------------------- */
/*                               COMANDO MKFILE                               */
/* -------------------------------------------------------------------------- */
func ProcessMKFILE(input string, path *string, r *bool, size *int, cont *string, flagN *bool) {
	padreBusqueda = 0
	Padre.B_inodo = -1
	copy(Padre.B_name[:], "")
	flags := strings.Split(input, "-")
	for _, i := range flags {
		if i == "r" {
			*r = true
		}
		f := strings.Split(i, "=")
		if f[0] == "path" {
			*path = f[1]

			if strings.Contains(f[1], " ") {
				tieneComillas := strings.Split(*path, "\"")
				if len(tieneComillas)-1 == 0 {
					*path = `"` + f[1] + `"`
				}

			}
		}

	}

	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
		case "r":
		case "size":
			sizeValue := 0
			fmt.Sscanf(flagValue, "%d", &sizeValue)
			*size = sizeValue
		case "cont":
			*cont = flagValue
		default:
		}
	}

	AddText("--------------------------------------------------------------------------")
	AddText("                   		MKFILE: PROCESANDO	 	                       ")
	AddText(fmt.Sprintf("\n" + *path + "\n"))
	AddText("--------------------------------------------------------------------------")

	re = regexp.MustCompile(`-(\w+)`)
	matches = re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]

		switch flagName {
		case "r":
			*r = true
		case "size":
		case "path":
		case "cont":
		default:
			AddText("Error: Flag not found: " + flagName)
			*flagN = true
		}
	}
}

func MKFILE(path *string, r *bool) {
	/* -------------------------------------------------------------------------- */
	/*                  COMPROBAMOS SI HAY UNA SESSION EXISTENTE                  */
	/* -------------------------------------------------------------------------- */
	if !session {
		AddText("--------------------------------------------------------------------------")
		AddText("                   MKFILE: NO HAY UNA SESION INICIADA                     ")
		AddText("--------------------------------------------------------------------------")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepaths := "./Disks/" + letra + ".dsk"
	file,err := utilities_test.OpenFile(filepaths)
	if err != nil {
		AddText("Error opening disk file:"+ err.Error())
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		AddText("Error reading MBR:"+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL INODO 0                            */
	/* -------------------------------------------------------------------------- */

	var Inode0 structs_test.Inode
	if err := utilities_test.ReadObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs_test.Inode{})))); err != nil {
		AddText("Error reading inode:"+ err.Error())
		return
	}


	/* -------------------------------------------------------------------------- */
	/*                           OBTENEMOS LA RUTA PADRE                          */
	/* -------------------------------------------------------------------------- */
	tieneComillas := strings.Split(*path, "\"")
	if len(tieneComillas)-1 != 0 {
		if len(tieneComillas)-1 == 1 {
			*path = tieneComillas[0]
		} else {
			*path = tieneComillas[1]
		}
	}
	rutaPadre := filepath.Dir(*path)
	AddText("Ruta original")
	AddText(*path)
	AddText("Ruta padre")
	AddText(rutaPadre)
	Carpetas := strings.Split(*path, "/")
	nuevaCarpeta := Carpetas[len(Carpetas)-1]
	partes := strings.Split(rutaPadre, "/")
	partes = partes[1:]
	//AddText("Elementos en ruta padre")
	//AddText(len(partes) - 1)
	carpetaCreada := false
	/* -------------------------------------------------------------------------- */
	/*                     RECORREMOS LOS BLOQUES DEL INODO 0                     */
	/* -------------------------------------------------------------------------- */
	//AddText("Bloques del inodo 0:")
	ultimo := 0
	root := false
	padreExiste := false
	for cont, i := range Inode0.I_block {
		if len(partes)-1 == 0 {
			//AddText("root es true")
			root = true
		}
		if i == -1 {
			ultimo = int(cont - 1)
			break
		}
		//AddText(i)
		if !root {
			existe := BuscarRuta(partes, i, 0)
			if existe {
				AddText("Existe la ruta padre")
				padreExiste = true
			}
		}
	}

	if root {
		existe := false
		for _, i := range Inode0.I_block {
			if i == -1 {
				break
			}
			// print("Buscando en el inodo 0 el bloque ")
			// AddText(i)
			existe = BuscarEspacioEnRoot(nuevaCarpeta, i)
			if existe {
				break
			}
		}
		if !existe {
			AddText("Creando nuevo inodo y bloque")
			// BlockCounter++
			CrrSuperblock.S_blocks_count++
			Inode0.I_block[ultimo+1] = CrrSuperblock.S_blocks_count
			CrearFolderBlock(file, CrrSuperblock.S_blocks_count, nuevaCarpeta)
			AddText("Actualizando inodo 0")
			if err := utilities_test.WriteObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs_test.Inode{})))); err != nil {
				AddText("Error reading inode:"+ err.Error())
				return
			}

		}
		carpetaCreada = true
	}

	if padreExiste && !carpetaCreada {
		AddText("Creando carpeta desde padre")
		CreandoCamino(Padre.B_inodo, nuevaCarpeta, file, partes)
		carpetaCreada = true
	}

	if *r && !carpetaCreada {
		if string(Padre.B_name[:]) != "" {
			AddText("creando a partir de carpetas existentes")
			AddText(fmt.Sprintf("Encontrado -> B_inode: %d B_name: %s\n", Padre.B_inodo, Padre.B_name))
			CreandoCamino(Padre.B_inodo, nuevaCarpeta, file, partes)
		} else {
			AddText("Creando todas las carpetas")
			CreandoCamino(0, nuevaCarpeta, file, partes)
		}
		carpetaCreada = true
	}
	if carpetaCreada {
		AddText("--------------------------------------------------------------------------")
		AddText(fmt.Sprintf("                MKDIR: CARPETA %s CREADA CORRECTAMENTE", nuevaCarpeta))
		AddText("--------------------------------------------------------------------------")
	} else {
		AddText("Error: No se logro crear la carpeta")
	}
}

/* -------------------------------------------------------------------------- */
/*                                 COMANDO CAT                                */
/* -------------------------------------------------------------------------- */
func ProcessCAT(input string, file *string) string {
	re := regexp.MustCompile(`-file(\d*)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagIndex := match[1]
		flagValue := match[2]

		// Eliminar comillas si están presentes en el valor
		flagValue = strings.Trim(flagValue, "\"")

		// Generar el nombre de la clave para el mapa

		// Asignar el valor al mapa
		*file = flagValue
		return flagIndex
	}
	return ""
}

func CAT(file *string) {
	/* -------------------------------------------------------------------------- */
	/*                  COMPROBAMOS SI HAY UNA SESSION EXISTENTE                  */
	/* -------------------------------------------------------------------------- */

	if !session {
		AddText("--------------------------------------------------------------------------")
		AddText("                   MKDIR: NO HAY UNA SESION INICIADA                      ")
		AddText("--------------------------------------------------------------------------")
		return
	}
}

/* -------------------------------------------------------------------------- */
/*                               COMANDO REMOVE                               */
/* -------------------------------------------------------------------------- */
func ProcessREMOVE(input string, path *string, flagN *bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		default:
			AddText("Error: Flag not found: " + flagName)
			*flagN = true
		}
	}
}

func REMOVE(path *string) {
	/* -------------------------------------------------------------------------- */
	/*                  COMPROBAMOS SI HAY UNA SESSION EXISTENTE                  */
	/* -------------------------------------------------------------------------- */

	if !session {
		AddText("--------------------------------------------------------------------------")
		AddText("                   MKDIR: NO HAY UNA SESION INICIADA                      ")
		AddText("--------------------------------------------------------------------------")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepaths := "./Disks/" + letra + ".dsk"
	file,err := utilities_test.OpenFile(filepaths)
	if err != nil {
		AddText("Error opening disk file:"+ err.Error())
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		AddText("Error reading MBR:"+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL INODO 0                            */
	/* -------------------------------------------------------------------------- */

	var Inode0 structs_test.Inode
	if err := utilities_test.ReadObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs_test.Inode{})))); err != nil {
		AddText("Error reading inode:"+ err.Error())
		return
	}

	
	carpetas := strings.Split(*path, "/")
	carpetas = carpetas[1:]
	nuevaCarpeta := carpetas[len(carpetas)-1]
	/* -------------------------------------------------------------------------- */
	/*                     RECORREMOS LOS BLOQUES DEL INODO 0                     */
	/* -------------------------------------------------------------------------- */
	//AddText("Bloques del inodo 0:")
	AddText("eliminando " + nuevaCarpeta)
	deleted := false
	for _, i := range Inode0.I_block {
		if i == -1 {
			break
		}
		//AddText(i)
		deleted = EliminarRuta(carpetas, i, 0)
		if deleted {
			AddText(nuevaCarpeta + " eliminado con exito")
			break
		}
	}
	if !deleted {
		AddText("No se logro eliminar " + nuevaCarpeta)
	}
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO EDIT                                */
/* -------------------------------------------------------------------------- */
func ProcessEDIT(input string, path *string, cont *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "cont":
			*cont = flagValue
		default:
			AddText("Error: Flag not found: " + flagName)
		}
	}
}

func EDIT(path *string, cont *string) {
	/* -------------------------------------------------------------------------- */
	/*                  COMPROBAMOS SI HAY UNA SESSION EXISTENTE                  */
	/* -------------------------------------------------------------------------- */

	if !session {
		AddText("--------------------------------------------------------------------------")
		AddText("                   MKDIR: NO HAY UNA SESION INICIADA                      ")
		AddText("--------------------------------------------------------------------------")
		return
	}
}

/* -------------------------------------------------------------------------- */
/*                               COMANDO RENAME                               */
/* -------------------------------------------------------------------------- */
func ProcessRENAME(input string, path *string, name *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "name":
			*name = flagValue
		default:
			AddText("Error: Flag not found: " + flagName)
		}
	}
}

func RENAME(path *string, name *string) {
	/* -------------------------------------------------------------------------- */
	/*                  COMPROBAMOS SI HAY UNA SESSION EXISTENTE                  */
	/* -------------------------------------------------------------------------- */

	if !session {
		AddText("--------------------------------------------------------------------------")
		AddText("                   MKDIR: NO HAY UNA SESION INICIADA                      ")
		AddText("--------------------------------------------------------------------------")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepaths := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepaths)
	if err != nil {
		AddText("Error opening disk file:"+ err.Error())
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		AddText("Error reading MBR:"+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL INODO 0                            */
	/* -------------------------------------------------------------------------- */

	var Inode0 structs_test.Inode
	if err := utilities_test.ReadObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs_test.Inode{})))); err != nil {
		AddText("Error reading inode:"+ err.Error())
		return
	}
	

	carpetas := strings.Split(*path, "/")
	carpetas = carpetas[1:]
	nuevaCarpeta := carpetas[len(carpetas)-1]
	/* -------------------------------------------------------------------------- */
	/*                     RECORREMOS LOS BLOQUES DEL INODO 0                     */
	/* -------------------------------------------------------------------------- */
	//AddText("Bloques del inodo 0:")
	for _, i := range Inode0.I_block {
		if i == -1 {
			break
		}
		//AddText(i)
		rename := Rename(carpetas, i, 0, *name)
		if rename {
			AddText(nuevaCarpeta + " renombrada con exito")
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO COPY                                */
/* -------------------------------------------------------------------------- */
func ProcessCOPY(input string, path *string, destino *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "destino":
			*destino = flagValue
		default:
			AddText("Error: Flag not found: " + flagName)
		}
	}
}

func COPY(path *string, destino *string) {
	/* -------------------------------------------------------------------------- */
	/*                  COMPROBAMOS SI HAY UNA SESSION EXISTENTE                  */
	/* -------------------------------------------------------------------------- */
	if !session {
		AddText("--------------------------------------------------------------------------")
		AddText("                   MKFILE: NO HAY UNA SESION INICIADA                     ")
		AddText("--------------------------------------------------------------------------")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepaths := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepaths)
	if err != nil {
		AddText("Error opening disk file:"+ err.Error())
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		AddText("Error reading MBR:"+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL INODO 0                            */
	/* -------------------------------------------------------------------------- */

	var Inode0 structs_test.Inode
	if err := utilities_test.ReadObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs_test.Inode{})))); err != nil {
		AddText("Error reading inode:"+ err.Error())
		return
	}


	/* -------------------------------------------------------------------------- */
	/*                           OBTENEMOS LA RUTA PADRE                          */
	/* -------------------------------------------------------------------------- */
	tieneComillas := strings.Split(*path, "\"")
	if len(tieneComillas)-1 != 0 {
		if len(tieneComillas)-1 == 1 {
			*path = tieneComillas[0]
		} else {
			*path = tieneComillas[1]
		}
	}
	tieneComillas = strings.Split(*destino, "\"")
	if len(tieneComillas)-1 != 0 {
		if len(tieneComillas)-1 == 1 {
			*destino = tieneComillas[0]
		} else {
			*destino = tieneComillas[1]
		}
	}
	Carpetas := strings.Split(*path, "/")
	nuevaCarpeta := Carpetas[len(Carpetas)-1]
	partes := strings.Split(*destino, "/")
	partes = partes[1:]
	//AddText("Elementos en ruta padre")
	//AddText(len(partes) - 1)
	carpetaCreada := false
	/* -------------------------------------------------------------------------- */
	/*                     RECORREMOS LOS BLOQUES DEL INODO 0                     */
	/* -------------------------------------------------------------------------- */
	//AddText("Bloques del inodo 0:")
	ultimo := 0
	root := false
	padreExiste := false
	for cont, i := range Inode0.I_block {
		if len(partes)-1 == 0 {
			//AddText("root es true")
			root = true
		}
		if i == -1 {
			ultimo = int(cont - 1)
			break
		}
		//AddText(i)
		if !root {
			existe := BuscarRuta(partes, i, 0)
			if existe {
				AddText("Existe la ruta padre")
				padreExiste = true
			}
		}
	}

	if root {
		existe := false
		for _, i := range Inode0.I_block {
			if i == -1 {
				break
			}
			// print("Buscando en el inodo 0 el bloque ")
			// AddText(i)
			existe = BuscarEspacioEnRoot(nuevaCarpeta, i)
			if existe {
				break
			}
		}
		if !existe {
			AddText("Creando nuevo inodo y bloque")
			// BlockCounter++
			CrrSuperblock.S_blocks_count++
			Inode0.I_block[ultimo+1] = CrrSuperblock.S_blocks_count
			CrearFolderBlock(file, CrrSuperblock.S_blocks_count, nuevaCarpeta)
			AddText("Actualizando inodo 0")
			if err := utilities_test.WriteObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs_test.Inode{})))); err != nil {
				AddText("Error reading inode:"+ err.Error())
				return
			}

		}
		carpetaCreada = true
	}

	if padreExiste && !carpetaCreada {
		AddText("Creando carpeta desde padre")
		CreandoCamino(Padre.B_inodo, nuevaCarpeta, file, partes)
		carpetaCreada = true
	}

	if !carpetaCreada {
		if string(Padre.B_name[:]) != "" {
			AddText("creando a partir de carpetas existentes")
			fmt.Printf("Encontrado -> B_inode: %d B_name: %s\n", Padre.B_inodo, Padre.B_name)
			CreandoCamino(Padre.B_inodo, nuevaCarpeta, file, partes)
		} else {
			AddText("Creando todas las carpetas")
			CreandoCamino(0, nuevaCarpeta, file, partes)
		}
		carpetaCreada = true
	}
	if carpetaCreada {
		AddText("--------------------------------------------------------------------------")
		AddText(fmt.Sprintf("                COPY:  %s COPIADO CORRECTAMENTE\n", nuevaCarpeta))
		AddText("--------------------------------------------------------------------------")
	} else {
		AddText("Error: No se logro copiar el elemento")
	}
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO MOVE                                */
/* -------------------------------------------------------------------------- */
func ProcessMOVE(input string, path *string, destino *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "destino":
			*destino = flagValue
		default:
			AddText("Error: Flag not found: " + flagName)
		}
	}
}

func MOVE(path *string, destino *string) {
	/* -------------------------------------------------------------------------- */
	/*                  COMPROBAMOS SI HAY UNA SESSION EXISTENTE                  */
	/* -------------------------------------------------------------------------- */

	if !session {
		AddText("--------------------------------------------------------------------------")
		AddText("                   MKDIR: NO HAY UNA SESION INICIADA                      ")
		AddText("--------------------------------------------------------------------------")
		return
	}

	REMOVE(path)
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepaths := "./Disks/" + letra + ".dsk"
	file,err := utilities_test.OpenFile(filepaths)
	if err != nil {
		AddText("Error opening disk file:"+ err.Error())
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		AddText("Error reading MBR:"+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL INODO 0                            */
	/* -------------------------------------------------------------------------- */

	var Inode0 structs_test.Inode
	if err := utilities_test.ReadObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs_test.Inode{})))); err != nil {
		AddText("Error reading inode:"+ err.Error())
		return
	}


	/* -------------------------------------------------------------------------- */
	/*                           OBTENEMOS LA RUTA PADRE                          */
	/* -------------------------------------------------------------------------- */
	tieneComillas := strings.Split(*path, "\"")
	if len(tieneComillas)-1 != 0 {
		if len(tieneComillas)-1 == 1 {
			*path = tieneComillas[0]
		} else {
			*path = tieneComillas[1]
		}
	}
	tieneComillas = strings.Split(*destino, "\"")
	if len(tieneComillas)-1 != 0 {
		if len(tieneComillas)-1 == 1 {
			*destino = tieneComillas[0]
		} else {
			*destino = tieneComillas[1]
		}
	}
	Carpetas := strings.Split(*path, "/")
	nuevaCarpeta := Carpetas[len(Carpetas)-1]
	partes := strings.Split(*destino, "/")
	partes = partes[1:]
	//AddText("Elementos en ruta padre")
	//AddText(len(partes) - 1)
	carpetaCreada := false
	/* -------------------------------------------------------------------------- */
	/*                     RECORREMOS LOS BLOQUES DEL INODO 0                     */
	/* -------------------------------------------------------------------------- */
	//AddText("Bloques del inodo 0:")
	ultimo := 0
	root := false
	padreExiste := false
	for cont, i := range Inode0.I_block {
		if len(partes)-1 == 0 {
			//AddText("root es true")
			root = true
		}
		if i == -1 {
			ultimo = int(cont - 1)
			break
		}
		//AddText(i)
		if !root {
			existe := BuscarRuta(partes, i, 0)
			if existe {
				AddText("Existe la ruta padre")
				padreExiste = true
			}
		}
	}

	if root {
		existe := false
		for _, i := range Inode0.I_block {
			if i == -1 {
				break
			}
			// print("Buscando en el inodo 0 el bloque ")
			// AddText(i)
			existe = BuscarEspacioEnRoot(nuevaCarpeta, i)

			if existe {
				break
			}
		}
		if !existe {
			AddText("Creando nuevo inodo y bloque")
			// BlockCounter++
			CrrSuperblock.S_blocks_count++
			Inode0.I_block[ultimo+1] = CrrSuperblock.S_blocks_count
			CrearFolderBlock(file, CrrSuperblock.S_blocks_count, nuevaCarpeta)
			AddText("Actualizando inodo 0")
			if err := utilities_test.WriteObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs_test.Inode{})))); err != nil {
				AddText("Error reading inode:"+ err.Error())
				return
			}

		}
		carpetaCreada = true
	}

	if padreExiste && !carpetaCreada {
		AddText("Creando carpeta desde padre")
		CreandoCamino(Padre.B_inodo, nuevaCarpeta, file, partes)
		carpetaCreada = true
	}

	if !carpetaCreada {
		if string(Padre.B_name[:]) != "" {
			AddText("creando a partir de carpetas existentes")
			AddText(fmt.Sprintf("Encontrado -> B_inode: %d B_name: %s\n", Padre.B_inodo, Padre.B_name))
			CreandoCamino(Padre.B_inodo, nuevaCarpeta, file, partes)
		} else {
			AddText("Creando todas las carpetas")
			CreandoCamino(0, nuevaCarpeta, file, partes)
		}
		carpetaCreada = true
	}
	if carpetaCreada {
		AddText("--------------------------------------------------------------------------")
		AddText(fmt.Sprintf("                MOVE:  %s  CORRECTAMENTE\n", nuevaCarpeta))
		AddText("--------------------------------------------------------------------------")
	} else {
		AddText("Error: No se logro mover el elemento")
	}
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO FIND                                */
/* -------------------------------------------------------------------------- */
func ProcessFIND(input string, path *string, destino *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "destino":
			*destino = flagValue
		default:
			AddText("Error: Flag not found: " + flagName)
		}
	}
}

func FIND(path *string, destino *string) {
	/* -------------------------------------------------------------------------- */
	/*                  COMPROBAMOS SI HAY UNA SESSION EXISTENTE                  */
	/* -------------------------------------------------------------------------- */

	if !session {
		AddText("--------------------------------------------------------------------------")
		AddText("                   MKDIR: NO HAY UNA SESION INICIADA                      ")
		AddText("--------------------------------------------------------------------------")
		return
	}
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO CHOWN                               */
/* -------------------------------------------------------------------------- */
func ProcessCHOWN(input string, path *string, user *string, r *bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "user":
			*user = flagValue
		case "r":
			*r = true
		default:
			AddText("Error: Flag not found: " + flagName)
		}
	}
}

func CHOWN(path *string, user *string, r *bool) {
	/* -------------------------------------------------------------------------- */
	/*                  COMPROBAMOS SI HAY UNA SESSION EXISTENTE                  */
	/* -------------------------------------------------------------------------- */

	if !session {
		AddText("--------------------------------------------------------------------------")
		AddText("                   MKDIR: NO HAY UNA SESION INICIADA                      ")
		AddText("--------------------------------------------------------------------------")
		return
	}
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO CHMOD                               */
/* -------------------------------------------------------------------------- */
func ProcessCHMOD(input string, path *string, ugo *string, r *bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "ugo":
			*ugo = flagValue
		case "r":
			*r = true
		default:
			AddText("Error: Flag not found: " + flagName)
		}
	}
}

func CHMOD(path *string, ugo *string, r *bool) {
	/* -------------------------------------------------------------------------- */
	/*                  COMPROBAMOS SI HAY UNA SESSION EXISTENTE                  */
	/* -------------------------------------------------------------------------- */

	if !session {
		AddText("--------------------------------------------------------------------------")
		AddText("                   MKDIR: NO HAY UNA SESION INICIADA                      ")
		AddText("--------------------------------------------------------------------------")
		return
	}
}
