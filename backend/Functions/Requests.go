package functions_test

import (
	"P2/Structs"
	"P2/Utilities"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

// ? --------------------------------------------------------------------------
// ?               FUNCION PARA ENVIAR LAS PARTICIONES DE UN DISCO
// ? --------------------------------------------------------------------------
func GetPartitions(driveletter string) ([]string, error) {
	var particiones []string
	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	rutaDisco := "./Disks/" + driveletter + ".dsk"
	file, err := os.Open(rutaDisco)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                               CARGAMOS EL MBR                              */
	/* -------------------------------------------------------------------------- */

	var TempMBR structs_test.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		return nil, err
	}

	/* -------------------------------------------------------------------------- */
	/*                         RECORREMOS LAS PARTICIONES                         */
	/* -------------------------------------------------------------------------- */
	for _, partition := range TempMBR.Mbr_particion {
		if string(partition.Part_status[:]) == "1" {
			partNameClean := strings.Trim(string(partition.Part_name[:]), "\x00")
			particiones = append(particiones, partNameClean)
		}
	}

	return particiones, nil
}

// ? --------------------------------------------------------------------------
// ?                FUNCION PARA RETORNAR EL RESULTADO DEL LOGIN
// ? --------------------------------------------------------------------------
func Session(driveletter string, name string, user string, password string) bool {
	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	rutaDisco := "./Disks/" + driveletter + ".dsk"
	file, err := os.Open(rutaDisco)
	if err != nil {
		return false
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                               CARGAMOS EL MBR                              */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		return false
	}

	/* -------------------------------------------------------------------------- */
	/*                         RECORREMOS LAS PARTICIONES                         */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		partNameClean := strings.Trim(string(TempMBR.Mbr_particion[i].Part_name[:]), "\x00")
		if TempMBR.Mbr_particion[i].Part_size != 0 && partNameClean == name {
			index = i
			break
		}
	}
	if index == -1 {
		return false
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		AddText("Error: " + err.Error())
		return false
	}

	/* -------------------------------------------------------------------------- */
	/*                  CARGAMOS EL ARCHIVO USERS.TXT DEL INODO 1                 */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		AddText("Error: " + err.Error())
		return false
	}

	/* -------------------------------------------------------------------------- */
	/*                              LEEMOS EL ARCHIVO                             */
	/* -------------------------------------------------------------------------- */
	limit := int32(0)
	for i := int32(0); i < 15; i++ {
		if crrInode.I_block[i] != -1 {
			limit++
		} else {
			break
		}
	}

	for i := int32(0); i < limit; i++ {
		var Fileblock structs_test.Fileblock
		blockNum := crrInode.I_block[i]
		if err := utilities_test.ReadObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
			AddText("Error: " + err.Error())
			return false
		}

		data := string(Fileblock.B_content[:])
		lines := strings.Split(data, "\n")

		for _, line := range lines {
			// Imprimir cada línea
			//fmt.Println(line)
			items := strings.Split(line, ",")
			if len(items) > 3 {
				//AddText("items[2]->" + items[2])
				if user == items[len(items)-2] {
					return true
				}
			}
		}
	}

	return false
}

// ?--------------------------------------------------------------------------
// ?             FUNCION PARA RETORNAR EL CONTENIDO DE UNA CARPETA
// ?--------------------------------------------------------------------------
func FolderContent(driveletter string, partition string, ruta string) []string {
	var contenido []string
	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	rutaDisco := "./Disks/" + driveletter + ".dsk"
	file, err := os.Open(rutaDisco)
	if err != nil {
		println("Error al leer el disco")
		return nil
	}
	defer file.Close()
	/* -------------------------------------------------------------------------- */
	/*                               CARGAMOS EL MBR                              */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		println("Error al cargar el mbr")
		return nil
	}

	/* -------------------------------------------------------------------------- */
	/*                            CARGAMOS LA PARTICION                           */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		partNameClean := strings.Trim(string(TempMBR.Mbr_particion[i].Part_name[:]), "\x00")
		if TempMBR.Mbr_particion[i].Part_size != 0 && partNameClean == partition {
			index = i
			break
		}
	}
	if index == -1 {
		println("No encontro la particion")
		return nil
	}
	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		println("error al cargar el superbloque")
		return nil
	}

	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL INODO 0                            */
	/* -------------------------------------------------------------------------- */
	var Inode0 structs_test.Inode
	if err := utilities_test.ReadObject(file, &Inode0, int64(tempSuperblock.S_inode_start+0*int32(binary.Size(structs_test.Inode{})))); err != nil {
		println("error al cargar el inodo 0")
		return nil
	}
	//structs_test.PrintInode(Inode0)
	/* -------------------------------------------------------------------------- */
	/*                             RECORREMOS LA RUTA                             */
	/* -------------------------------------------------------------------------- */
	carpetas := strings.Split(ruta, "/")
	partes := carpetas[1:]
	if len(carpetas[1]) == 0 {
		for cont, i := range Inode0.I_block {
			if i == -1 {
				break
			}

			bloque := CargarBloque(file, Inode0.I_block[cont])
			for _, j := range bloque.B_content {
				partNameClean := strings.Trim(string(j.B_name[:]), "\x00")
				if partNameClean == "." {
					continue
				}
				if partNameClean == ".." {
					continue
				}
				if len(partNameClean) == 0 {
					continue
				}
				contenido = append(contenido, partNameClean)
			}
		}
		return contenido
	} else {
		for _, i := range Inode0.I_block {
			if i == -1 {
				break
			}

			contenido = BuscarContenido(file, partes, i, 0, tempSuperblock)
			if contenido != nil {
				return contenido
			}
		}
	}
	return nil
}

func BuscarContenido(file *os.File, ruta []string, bloque int32, busqueda int, sp structs_test.Superblock) []string {
	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL BLOQUE                             */
	/* -------------------------------------------------------------------------- */
	var FolderBlock structs_test.Folderblock
	//var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &FolderBlock, int64(sp.S_block_start+bloque*int32(binary.Size(structs_test.Folderblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return nil
	}
	encontrado := false
	for _, content := range FolderBlock.B_content {
		nombre := string(bytes.Trim(content.B_name[:], "\x00")) // Eliminar bytes nulos del final
		if strings.TrimSpace(nombre) == strings.TrimSpace(ruta[busqueda]) {
			encontrado = true
			//Cargar el inodo
			//recorrerlos bloques del inodo
			//recursividad para cada bloque hasta que se encuentre la otra parte de la ruta
			Inode := CargarInodo(file, content.B_inodo)
			if busqueda < len(ruta)-1 {
				busqueda++
				for _, i := range Inode.I_block {
					if i == -1 {
						break
					}
					resultado := BuscarContenido(file, ruta, i, busqueda, sp)
					if resultado != nil{
						return resultado
					}

				}
			} else {
				var contenido []string
				for cont, i := range Inode.I_block {
					if i == -1 {
						return contenido
					}

					bloque := CargarBloque(file, Inode.I_block[cont])
					for _, j := range bloque.B_content {
						partNameClean := strings.Trim(string(j.B_name[:]), "\x00")
						if partNameClean == "." {
							continue
						}
						if partNameClean == ".." {
							continue
						}
						if len(partNameClean) == 0 {
							break
						}
						contenido = append(contenido, partNameClean)
					}
				}
			}
		}
		if encontrado {
			break
		}
	}

	return nil
}
