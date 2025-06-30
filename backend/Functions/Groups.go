package functions_test

import (
	"P2/Global"
	"P2/Structs"
	"P2/Utilities"
	"encoding/binary"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	session       = false
	usuario       = Global.UserInfo{}
	groupCounter  = 1
	userCounter   = 1
	InodeIndex    = int32(1)
	blockIndex    = 0
	searchIndex   = 0
	letra         = ""
	ID            = ""
	CrrSuperblock structs_test.Superblock
	indexSB       = 0
)

//?                    ADMINISTRACION DE USUARIOS Y GRUPOS
/* -------------------------------------------------------------------------- */
/*                                COMANDO LOGIN                               */
/* -------------------------------------------------------------------------- */
func ProcessLOGIN(input string, user *string, pass *string, id *string, flagN *bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "user":
			*user = flagValue
		case "pass":
			*pass = flagValue
		case "id":
			*id = flagValue
		default:
			AddText("Error: Flag not found: " + flagName)
			*flagN = true
		}
	}
}

func LOGIN(user *string, pass *string, id *string) {

	letra = string((*id)[0])
	ID = *id

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		AddText("Error: "+ err.Error())
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		AddText("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	indexSB = index
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	CrrSuperblock = tempSuperblock

	/* -------------------------------------------------------------------------- */
	/*                   LEEMOS EL INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                             LEEMOS EL FILEBLOCK                            */
	/* -------------------------------------------------------------------------- */
	var Fileblock structs_test.Fileblock
	blockNum := crrInode.I_block[searchIndex]
	
	if err := utilities_test.ReadObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		AddText("Error: "+ err.Error())
		return
	}
	data := string(Fileblock.B_content[:])
	lines := strings.Split(data, "\n")

	userFound := false
	for _, line := range lines {
		// Imprimir cada línea
		// AddText(line)
		items := strings.Split(line, ",")
		if len(items) > 3 {
			//AddText("items[2]->" + items[2])
			if *user == items[len(items)-2] {
				userFound = true
				usuario.Nombre = items[len(items)-2]
				identificacion, err := strconv.Atoi(items[0])
				if err != nil {
					AddText("Error: "+ err.Error())
					return
				}
				usuario.ID = int32(identificacion)
				session = true
				break
			}
		}
	}

	if !userFound {
		searchIndex++
		if searchIndex <= blockIndex {
			LOGIN(user, pass, id)
		} else {
			AddText("Error: no se encontro al usuario")
			searchIndex = 0
			return
		}
	} else {
		//Global.PrintUser(usuario)
		searchIndex = 0
		return
	}
}

/* -------------------------------------------------------------------------- */
/*                               COMANDO LOGOUT                               */
/* -------------------------------------------------------------------------- */
func LOGOUT() {
	if session {
		AddText("--------------------------------------------------------------------------")
		AddText("                        LOGOUT: SESION CERRADA                            ")
		AddText("--------------------------------------------------------------------------")
		session = false
		searchIndex = 0
		usuario.Nombre = ""
		usuario.ID = -1
		return
	}
	println("Error: no hay una sesion activa")
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO MKGRP                               */
/* -------------------------------------------------------------------------- */
func ProcessMKGRP(input string, name *string, flagN *bool) {
	if usuario.Nombre == "root" {
		re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

		matches := re.FindAllStringSubmatch(input, -1)

		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			// Delete quotes if they are present in the value
			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "name":
				*name = flagValue
			default:
				AddText("Error: Flag not found: " + flagName)
				*flagN = true
			}
		}
	} else {
		println("Error: solo el usuario root puede acceder a este comando")
		return
	}
}

func MKGRP(name *string) {
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		AddText("Error: "+ err.Error())
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		AddText("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                   LEEMOS EL INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                             LEEMOS EL FILEBLOCK                            */
	/* -------------------------------------------------------------------------- */
	var Fileblock structs_test.Fileblock
	blockNum := crrInode.I_block[blockIndex]

	// if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
	// 	AddText("Error reading Fileblock:", err)
	// 	return
	// }
	if err := utilities_test.ReadObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	data := string(Fileblock.B_content[:])
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	/* -------------------------------------------------------------------------- */
	/*          ITERAMOS EN CADA LINEA PARA QUE NO HAYAN GRUPOS REPETIDOS         */
	/* -------------------------------------------------------------------------- */
	for _, line := range lines {
		// Imprimir cada línea
		AddText(line)
		items := strings.Split(line, ",")
		if len(items) == 3 {
			if *name == items[2] {
				AddText("Error: nombre repetido")
				return
			}
		}
	}

	/* -------------------------------------------------------------------------- */
	/*                          PARSEAMOS LA INFORMACION                          */
	/* -------------------------------------------------------------------------- */
	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
	groupCounter++
	nuevoGrupo := fmt.Sprintf("%d,G,%s\n", groupCounter, *name)
	newContent := currentContent + nuevoGrupo

	/* -------------------------------------------------------------------------- */
	/*                 CREAMOS MAS FILEBLOCKS PARA GUARDAR LA INFO                */
	/* -------------------------------------------------------------------------- */
	if len(newContent) > len(Fileblock.B_content) {
		if blockIndex > int(len(crrInode.I_block)) {
			AddText("Error: no hay mas bloques disponibles")
			return
		}
		blockIndex++
		//BlockCounter++
		CrrSuperblock.S_blocks_count++

		var NEWFileblock structs_test.Fileblock
		// if err := utilities_test.WriteObject(file, &NEWFileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
		// 	AddText("Error reading Fileblock:", err)
		// 	return
		// }
		if err := utilities_test.WriteObject(file, &NEWFileblock, int64(CrrSuperblock.S_block_start+CrrSuperblock.S_blocks_count*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
			AddText("Error: "+ err.Error())
			return
		}

		/* -------------------------------------------------------------------------- */
		/*                     ACTUALIZAMOS LOS BLOQUES DEL INODO 1                   */
		/* -------------------------------------------------------------------------- */
		crrInode.I_block[blockIndex] = CrrSuperblock.S_blocks_count
		if err := utilities_test.WriteObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
			AddText("Error: "+ err.Error())
			return
		}
		/* -------------------------------------------------------------------------- */
		/*                         ACTUALIZAMOS EL SUPERBLOQUE                        */
		/* -------------------------------------------------------------------------- */
		if err := utilities_test.WriteObject(file, &CrrSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
			AddText("Error: "+ err.Error())
			return
		}
		MKGRP(name)
	} else {
		/* -------------------------------------------------------------------------- */
		/*                GUARDA LA INFORMACION EN EL FILEBLOCK ACTUAL                */
		/* -------------------------------------------------------------------------- */
		copy(Fileblock.B_content[:], newContent)

		// if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
		// 	AddText("Error writing Fileblock to disk:", err)
		// 	return
		// }
		blockNum := crrInode.I_block[blockIndex]

		if err := utilities_test.WriteObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
			AddText("Error: "+ err.Error())
			return
		}

		//println("ACTUALIZACION")
		// Mostrar el contenido actualizado del Fileblock
		data := string(Fileblock.B_content[:])
		// Dividir la cadena en líneas
		lines := strings.Split(data, "\n")

		/* -------------------------------------------------------------------------- */
		/*          ITERAMOS EN CADA LINEA PARA QUE NO HAYAN GRUPOS REPETIDOS         */
		/* -------------------------------------------------------------------------- */
		for _, line := range lines {
			// Imprimir cada línea
			AddText(line)
		}
		/* -------------------------------------------------------------------------- */
		/*                         ACTUALIZAMOS EL SUPERBLOQUE                        */
		/* -------------------------------------------------------------------------- */
		if err := utilities_test.WriteObject(file, &CrrSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
			AddText("Error: "+ err.Error())
			return
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO RMGRP                               */
/* -------------------------------------------------------------------------- */
func ProcessRMGRP(input string, name *string, flagN *bool) {
	if usuario.Nombre == "root" {
		re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

		matches := re.FindAllStringSubmatch(input, -1)

		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			// Delete quotes if they are present in the value
			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "name":
				*name = flagValue
			default:
				AddText("Error: Flag not found: " + flagName)
				*flagN = true
			}
		}
	} else {
		println("Error: solo el usuario root puede acceder a este comando")
		return
	}
}

func RMGRP(name *string) {
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		AddText("Error: "+ err.Error())
		return
	}
	defer file.Close()

	// Leer el MBR del disco
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                               CARGAMOS EL MBR                              */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		AddText("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                   LEEMOS EN INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                      LEEMOS EL CONTENIDO DEL FILEBLOCK                     */
	/* -------------------------------------------------------------------------- */
	var Fileblock structs_test.Fileblock
	blockNum := crrInode.I_block[searchIndex]

	// if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
	// 	AddText("Error reading Fileblock:", err)
	// 	return
	// }
	if err := utilities_test.ReadObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                      COLOCAMOS EL STATUS DE ELIMINADO                      */
	/* -------------------------------------------------------------------------- */
	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
	lines := strings.Split(currentContent, "\n")
	deleted := false
	for i, line := range lines {
		if strings.Contains(line, *name) {
			lines[i] = "0,G," + *name
			deleted = true
			break
		}
	}

	/* -------------------------------------------------------------------------- */
	/*                   VERIFICAMOS BLOQUES O MENSAJE NOT FOUND                  */
	/* -------------------------------------------------------------------------- */
	if !deleted {
		searchIndex++
		if searchIndex > blockIndex {
			AddText("Group not found")
			searchIndex = 0
			return
		}
		RMGRP(name)

	}

	/* -------------------------------------------------------------------------- */
	/*                          ACTUALIZAMOS EL CONTENIDO                         */
	/* -------------------------------------------------------------------------- */
	newContent := strings.Join(lines, "\n")
	copy(Fileblock.B_content[:], newContent)

	if deleted {
		blockNum := crrInode.I_block[searchIndex]

		// if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
		// 	AddText("Error writing Fileblock to disk:", err)
		// 	return
		// }

		if err := utilities_test.WriteObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
			AddText("Error: "+ err.Error())
			return
		}

		currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
		lines := strings.Split(currentContent, "\n")
		for i := range lines {
			println(lines[i])
		}

		searchIndex = 0
	}
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO MKUSR                               */
/* -------------------------------------------------------------------------- */
func ProcessMKUSR(input string, user *string, pass *string, grp *string, flagN *bool) {
	if usuario.Nombre == "root" {
		re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

		matches := re.FindAllStringSubmatch(input, -1)

		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			// Delete quotes if they are present in the value
			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "user":
				*user = flagValue
			case "pass":
				*pass = flagValue
			case "grp":
				*grp = flagValue
			default:
				AddText("Error: Flag not found: " + flagName)
				*flagN = true
			}
		}
	} else {
		println("Error: solo el usuario root puede acceder a este comando")
		return
	}
}

func MKUSR(user *string, pass *string, grp *string) {
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		AddText("Error: "+ err.Error())
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		AddText("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                   LEEMOS EL INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	// AddText("Bitmap de bloques del inodo1")
	// AddText(crrInode.I_block)

	/* -------------------------------------------------------------------------- */
	/*                             LEEMOS EL FILEBLOCK                            */
	/* -------------------------------------------------------------------------- */
	blockNum := crrInode.I_block[blockIndex]
	var Fileblock structs_test.Fileblock
	// if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
	// 	AddText("Error reading Fileblock:", err)
	// 	return
	// }
	if err := utilities_test.ReadObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                          PARSEAMOS LA INFORMACION                          */
	/* -------------------------------------------------------------------------- */
	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
	groupCounter++
	searchIndex = 0
	var nuevoUsuario = BuscarGrupo(user, pass, grp)
	//AddText("nuevo usuarios: " + nuevoUsuario)
	if nuevoUsuario == "" {
		AddText("Error: No se encontro el grupo")
		return
	}
	newContent := currentContent + nuevoUsuario

	/* -------------------------------------------------------------------------- */
	/*                 CREAMOS MAS FILEBLOCKS PARA GUARDAR LA INFO                */
	/* -------------------------------------------------------------------------- */
	if len(newContent) > len(Fileblock.B_content) {
		if blockIndex > int(len(crrInode.I_block)) {
			AddText("Error: no hay mas bloques disponibles")
			return
		}
		blockIndex++
		//BlockCounter++
		CrrSuperblock.S_blocks_count++

		var NEWFileblock structs_test.Fileblock
		copy(NEWFileblock.B_content[:], nuevoUsuario)
		// if err := utilities_test.WriteObject(file, &NEWFileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
		// 	AddText("Error reading Fileblock:", err)
		// 	return
		// }

		if err := utilities_test.WriteObject(file, &NEWFileblock, int64(CrrSuperblock.S_block_start+CrrSuperblock.S_blocks_count*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
			AddText("Error: "+ err.Error())
			return
		}
		AddText("MKUSR EXITOSO")
		// Mostrar el contenido actualizado del Fileblock
		data := string(NEWFileblock.B_content[:])
		// Dividir la cadena en líneas
		lines := strings.Split(data, "\n")

		for _, line := range lines {
			// Imprimir cada línea
			AddText(line)
		}

		/* -------------------------------------------------------------------------- */
		/*                     ACTUALIZAMOS LOS BLOQUES DEL INODO                     */
		/* -------------------------------------------------------------------------- */
		crrInode.I_block[blockIndex] = CrrSuperblock.S_blocks_count

		if err := utilities_test.WriteObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
			AddText("Error: "+ err.Error())
			return
		}
		searchIndex = 0

		/* -------------------------------------------------------------------------- */
		/*                         ACTUALIZAMOS EL SUPERBLOQUE                        */
		/* -------------------------------------------------------------------------- */
		if err := utilities_test.WriteObject(file, &CrrSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
			AddText("Error: "+ err.Error())
			return
		}

	} else {
		println("MKUSR EXITOSO")
		/* -------------------------------------------------------------------------- */
		/*                GUARDA LA INFORMACION EN EL FILEBLOCK ACTUAL                */
		/* -------------------------------------------------------------------------- */
		copy(Fileblock.B_content[:], newContent)

		// if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
		// 	AddText("Error writing Fileblock to disk:", err)
		// 	return
		// }

		blockNum := crrInode.I_block[blockIndex]

		if err := utilities_test.WriteObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
			AddText("Error: "+ err.Error())
			return
		}

		// Mostrar el contenido actualizado del Fileblock
		data := string(Fileblock.B_content[:])
		// Dividir la cadena en líneas
		lines := strings.Split(data, "\n")

		/* -------------------------------------------------------------------------- */
		/*          ITERAMOS EN CADA LINEA PARA QUE NO HAYAN GRUPOS REPETIDOS         */
		/* -------------------------------------------------------------------------- */
		for _, line := range lines {
			// Imprimir cada línea
			AddText(line)
		}
		searchIndex = 0

		/* -------------------------------------------------------------------------- */
		/*                         ACTUALIZAMOS EL SUPERBLOQUE                        */
		/* -------------------------------------------------------------------------- */
		if err := utilities_test.ReadObject(file, &CrrSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
			AddText("Error: "+ err.Error())
			return
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO RMUSR                               */
/* -------------------------------------------------------------------------- */
func ProcessRMUSR(input string, user *string, flagN *bool) {
	if usuario.Nombre == "root" {
		re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

		matches := re.FindAllStringSubmatch(input, -1)

		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			// Delete quotes if they are present in the value
			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "user":
				*user = flagValue
			default:
				AddText("Error: Flag not found: " + flagName)
				*flagN = true
			}
		}
	} else {
		println("Error: solo el usuario root puede acceder a este comando")
		return
	}
}

func RMUSR(user *string) {
	filepath := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		AddText("Error: "+ err.Error())
		return
	}
	defer file.Close()

	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		AddText("Partition not found")
		return
	}

	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	var Fileblock structs_test.Fileblock
	// if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
	// 	AddText("Error reading Fileblock:", err)
	// 	return
	// }
	blockNum := crrInode.I_block[searchIndex]

	if err := utilities_test.ReadObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	data := string(Fileblock.B_content[:])
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		items := strings.Split(line, ",")
		if len(items) > 3 {
			if *user == items[len(items)-2] {
				items[0] = "0" // Setear el ID a 0
				newLine := strings.Join(items, ",")
				copy(Fileblock.B_content[:], []byte(newLine))
				// if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
				// 	AddText("Error writing Fileblock to disk:", err)
				// 	return
				// }
				blockNum := crrInode.I_block[searchIndex]

				if err := utilities_test.WriteObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
					AddText("Error: "+ err.Error())
					return
				}
				println("RMUSR " + *user + " exitoso")
				return
			}
		}
	}

	searchIndex++
	if searchIndex <= blockIndex {
		RMUSR(user)
	} else {
		AddText("User not found")
	}
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO CHGRP                               */
/* -------------------------------------------------------------------------- */
func ProcessCHGRP(input string, user *string, grp *string, flagN *bool) {
	if usuario.Nombre == "root" {
		searchIndex = 0
		re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

		matches := re.FindAllStringSubmatch(input, -1)

		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			// Delete quotes if they are present in the value
			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "user":
				*user = flagValue
			case "grp":
				*grp = flagValue
			default:
				AddText("Error: Flag not found: " + flagName)
				*flagN = true
			}
		}
	} else {
		println("Error: solo el usuario root puede acceder a este comando")
		return
	}
}

func CHGRP(user *string, grp *string) {
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		AddText("Error: "+ err.Error())
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		AddText("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	indexSB = index
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	CrrSuperblock = tempSuperblock

	/* -------------------------------------------------------------------------- */
	/*                   LEEMOS EL INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		AddText("Error: "+ err.Error())
		return
	}

	// AddText("Bitmap de bloques del inodo1")
	// AddText(crrInode.I_block)

	/* -------------------------------------------------------------------------- */
	/*                             LEEMOS EL FILEBLOCK                            */
	/* -------------------------------------------------------------------------- */
	var Fileblock structs_test.Fileblock
	blockNum := crrInode.I_block[searchIndex]
	// if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
	// 	AddText("Error reading Fileblock:", err)
	// 	return
	// }
	if err := utilities_test.ReadObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		AddText("Error: "+ err.Error())
		return
	}
	//AddText("Fileblock " + fmt.Sprint(searchIndex))
	data := string(Fileblock.B_content[:])
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	userFound := false
	for _, line := range lines {
		// Imprimir cada línea
		//AddText(line)
		items := strings.Split(line, ",")
		if len(items) > 3 {
			//AddText("items[2]->" + items[2])
			if *user == items[len(items)-2] {
				//print(items[2])
				items[2] = *grp // cambiar el grupo
				newLine := strings.Join(items, ",")
				copy(Fileblock.B_content[:], []byte(newLine))
				// if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
				// 	AddText("Error writing Fileblock to disk:", err)
				// 	return
				// }
				blockNum := crrInode.I_block[searchIndex]

				if err := utilities_test.WriteObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
					AddText("Error: "+ err.Error())
					return
				}
				println("RMUSR " + *user + " exitoso")
				return
			}
		}
	}

	if !userFound {
		searchIndex++
		if searchIndex <= blockIndex {
			CHGRP(user, grp)
		} else {
			AddText("Error: no se encontro al usuario")
			searchIndex = 0
			return
		}
	} else {
		Global.PrintUser(usuario)
		searchIndex = 0
		return
	}
}

/* -------------------------------------------------------------------------- */
/*                                 AUXILIARES                                 */
/* -------------------------------------------------------------------------- */

func BuscarGrupo(user *string, pass *string, grp *string) string {
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		AddText("Error: "+ err.Error())
		return ""
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		AddText("Error: "+ err.Error())
		return ""
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		AddText("Partition not found")
		return ""
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		AddText("Error: "+ err.Error())
		return ""
	}

	/* -------------------------------------------------------------------------- */
	/*                   LEEMOS EL INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		AddText("Error: "+ err.Error())
		return ""
	}

	// AddText("Bitmap de bloques del inodo1")
	// AddText(crrInode.I_block)

	/* -------------------------------------------------------------------------- */
	/*                             LEEMOS EL FILEBLOCK                            */
	/* -------------------------------------------------------------------------- */
	var Fileblock structs_test.Fileblock
	blockNum := crrInode.I_block[searchIndex]

	// if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
	// 	AddText("Error reading Fileblock:", err)
	// 	return ""
	// }
	if err := utilities_test.ReadObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		AddText("Error: "+ err.Error())
		return ""
	}
	//AddText("Fileblock " + fmt.Sprint(searchIndex))
	data := string(Fileblock.B_content[:])
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	groupFound := false
	var newUserLine string
	for _, line := range lines {
		// Imprimir cada línea
		//AddText(line)
		items := strings.Split(line, ",")
		if len(items) == 3 {
			//AddText("items[2]->" + items[2])
			if *grp == items[2] {
				groupFound = true
				newUserLine = fmt.Sprintf("%d,G,%s,%s,%s\n", userCounter, *grp, *user, *pass)
				userCounter++
				break
			}
		}
	}

	if !groupFound {
		searchIndex++
		if searchIndex <= blockIndex {
			return BuscarGrupo(user, pass, grp)
		}
	} else {
		return newUserLine
	}
	return ""
}
