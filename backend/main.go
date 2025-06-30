package main

import (
	analyzer_test "P2/Analyzer"
	functions_test "P2/Functions"
	utilities_test "P2/Utilities"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var arranque = true

func EmptyFolder(carpeta string) {
	archivos, err := filepath.Glob(filepath.Join(carpeta, "*"))
	if err != nil {
		fmt.Println("Error al obtener la lista de archivos:", err)
		return
	}

	for _, archivo := range archivos {
		err := os.Remove(archivo)
		if err != nil {
			fmt.Printf("Error al eliminar el archivo %s: %s\n", archivo, err)
		}
	}

	archivos, err = filepath.Glob(filepath.Join(carpeta, "*"))
	if err != nil {
		fmt.Println("Error al obtener la lista de archivos:", err)
		return
	}

	if len(archivos) != 0 {
		fmt.Println("La carpeta no pudo ser vaciada completamente.")
	}
}

func FolderInfo(carpeta string) ([]string, error) {
	var nombres []string
	archivos, err := os.ReadDir(carpeta)
	if err != nil {
		return nil, err
	}
	for _, archivo := range archivos {
		if !archivo.IsDir() {
			info, err := archivo.Info() // Accede a la información del archivo
			if err != nil {
				return nil, err
			}
			nombres = append(nombres, info.Name())
		}
	}
	return nombres, nil
}

func main() {

	app := fiber.New()
	// Solucion a CORS
	app.Use(cors.New())

	/* -------------------------------------------------------------------------- */
	/*                               MENSAJE INICIAL                              */
	/* -------------------------------------------------------------------------- */
	app.Get("/", func(c *fiber.Ctx) error {
		if arranque {
			//Vaciamos las carpetas
			EmptyFolder("./Disks")
			EmptyFolder("./Reports")
			arranque = false
		}
		return c.SendString("Ingresa un comando")
	})

	/* -------------------------------------------------------------------------- */
	/*                            COMANDOS INDIVIDUALES                           */
	/* -------------------------------------------------------------------------- */
	app.Post("/command", func(c *fiber.Ctx) error {
		// Obtener el texto enviado en el cuerpo de la solicitud
		text := c.FormValue("texto")
		//fmt.Println(text)
		analyzer_test.Command(text)
		// Enviar respuesta de éxito al cliente
		return c.SendString(utilities_test.Resultados.String())
	})

	/* -------------------------------------------------------------------------- */
	/*                             MANEJO DE ARCHIVOS                             */
	/* -------------------------------------------------------------------------- */
	app.Post("/upload", func(c *fiber.Ctx) error {
		// Recibe el contenido como texto
		body := c.FormValue("fileContent")

		// Divide el texto en líneas
		lines := strings.Split(body, "\n")

		// Procesa cada línea con la función analizar
		for _, line := range lines {
			analyzer_test.Command(line)
		}

		return c.SendString(utilities_test.Resultados.String())
	})

	/* -------------------------------------------------------------------------- */
	/*                         ENVIO DE DISCOS EXISTENTES                         */
	/* -------------------------------------------------------------------------- */
	app.Get("/disks", func(c *fiber.Ctx) error {
		carpeta := "./Disks" // Asegúrate de cambiar esto por la ruta correcta
		archivos, err := FolderInfo(carpeta)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(archivos) // Envía la lista de archivos como JSON
	})

	/* -------------------------------------------------------------------------- */
	/*                              ENVIO DE REPORTES                             */
	/* -------------------------------------------------------------------------- */
	app.Get("/reports", func(c *fiber.Ctx) error {
		archivos := utilities_test.GraphReports

		return c.JSON(archivos) // Envía la lista filtrada de archivos como JSON
	})

	/* -------------------------------------------------------------------------- */
	/*                            ENVIO DE PARTICIONES                            */
	/* -------------------------------------------------------------------------- */
	app.Post("/partitions", func(c *fiber.Ctx) error {
		body := c.FormValue("letter")
		archivos, err := functions_test.GetPartitions(body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.JSON(archivos)
	})

	/* -------------------------------------------------------------------------- */
	/*                               ENVIAR REPORTES                              */
	/* -------------------------------------------------------------------------- */
	app.Post("/reports/show", func(c *fiber.Ctx) error {
		type request struct {
			Filename string `json:"filename"`
		}

		var body request
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Error parsing JSON")
		}

		filePath := "./Reports/" + body.Filename // Asegúrate de validar y sanear la entrada para prevenir path traversal
		return c.Download(filePath)              // Envía el archivo como una descarga
	})

	/* -------------------------------------------------------------------------- */
	/*                                    LOGIN                                   */
	/* -------------------------------------------------------------------------- */

	app.Post("/login", func(c *fiber.Ctx) error {
		type LoginRequest struct {
			Username  string `json:"username"`
			Password  string `json:"password"`
			Disk      string `json:"disk"`
			Partition string `json:"partition"`
		}

		var request LoginRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false})
		}

		isAuthorized := functions_test.Session(request.Disk, request.Partition, request.Username, request.Password)

		// Enviar la respuesta
		return c.JSON(fiber.Map{"success": isAuthorized})
	})

	/* -------------------------------------------------------------------------- */
	/*                          CARGAR CONTENIDO DE RUTA                          */
	/* -------------------------------------------------------------------------- */
	app.Post("/docs", func(c *fiber.Ctx) error {
		type DocsRequest struct {
			Disk      string `json:"disk"`
			Partition string `json:"partition"`
			Ruta      string `json:"ruta"`
		}

		var request DocsRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Error parsing request body")
		}

		archivos := functions_test.FolderContent(request.Disk, request.Partition, request.Ruta)

		return c.JSON(archivos)
	})

	log.Fatal(app.Listen("0.0.0.0:4000"))
}
