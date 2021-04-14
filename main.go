package main

import (
	"bufio"        // Lecturade lineas incluyendo espacios
	"database/sql" // Inte con las bases de datos
	"fmt"
	"os"      // Bufer de lecturaos.Stdin
	"strconv" //Para convertir String a tipos de datos basicos

	_ "github.com/go-sql-driver/mysql" // La librería que nos permite conectar a MySQL
)

type ticket struct {
	Id                  int    `json:"Id"`
	User                string `json:"User"`
	Fecha_creacion      string `json:"Fecha_creacion"`
	Fecha_actualizacion string `json:"Fecha_actualizacion"`
	Estatus             bool   `json:"Estatus"`
}

func obtenerBaseDeDatos() (db *sql.DB, e error) {
	usuario := "root"
	pass := "1901"
	host := "tcp(localhost)"
	nombreBaseDeDatos := "tickets"
	// Debe tener la forma usuario:contraseña@host/nombreBaseDeDatos
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", usuario, pass, host, nombreBaseDeDatos))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func insertar(c ticket) (e error) {
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return err
	}
	defer db.Close()

	sentenciaPreparada, err := db.Prepare("INSERT INTO tickets_guardados (id_ticket, usuario, fecha_creacion, fecha_actualizacion, estatus) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer sentenciaPreparada.Close()

	_, err = sentenciaPreparada.Exec(c.Id, c.User, c.Fecha_creacion, c.Fecha_actualizacion, c.Estatus)
	if err != nil {
		return err
	}
	return nil
}

func obtenerTickets() ([]ticket, error) {
	var tickets []ticket
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	filas, err := db.Query("SELECT id_ticket, usuario, fecha_creacion, fecha_actualizacion, estatus FROM tickets_guardados")

	if err != nil {
		return nil, err
	}

	defer filas.Close() // en el caso de cerrar porque hay ejecución sin error

	// Mapeo de los datos traidos
	var t ticket

	for filas.Next() {
		err = filas.Scan(&t.Id, &t.User, &t.Fecha_creacion, &t.Fecha_actualizacion, &t.Estatus)
		// Prevenir el error al escanear
		if err != nil {
			return nil, err
		}

		tickets = append(tickets, t)
	}

	return tickets, nil
}

func actualizar(t ticket) error {
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return err
	}
	defer db.Close()

	sentenciaPreparada, err := db.Prepare("UPDATE tickets_guardados SET usuario = ?, fecha_creacion = ?, fecha_actualizacion = ?, estatus=? WHERE id_ticket = ?")
	if err != nil {
		return err
	}
	defer sentenciaPreparada.Close()
	_, err = sentenciaPreparada.Exec(t.User, t.Fecha_creacion, t.Fecha_actualizacion, t.Estatus, t.Id)
	return err
}

func eliminar(t ticket) error {
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return err
	}
	defer db.Close()

	sentenciaPreparada, err := db.Prepare("DELETE FROM tickets_guardados WHERE id_ticket = ?")
	if err != nil {
		return err
	}
	defer sentenciaPreparada.Close()

	_, err = sentenciaPreparada.Exec(t.Id)
	if err != nil {
		return err
	}
	return nil
}

func filtrar(user string) (ticket, error) {
	var t ticket
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return t, err
	}
	defer db.Close()

	defer db.Close()

	sentenciaPreparada, err := db.Query("SELECT * FROM tickets_guardados WHERE usuario = ?", user)
	if err != nil {
		return t, err
	}
	defer sentenciaPreparada.Close()
	for sentenciaPreparada.Next() {
		err = sentenciaPreparada.Scan(&t.User, &t.Fecha_creacion, &t.Fecha_actualizacion, &t.Estatus, &t.Id)
		if err != nil {
			return t, err
		}
	}
	return t, err
}

func main() {
	menu := `¿Qué deseas hacer?
			[1] -- Insertar
			[2] -- Mostrar
			[3] -- Filtrar
			[4] -- Actualizar
			[5] -- Eliminar
			[6] -- Salir
			----->	`
	var eleccion int
	var t ticket
	for eleccion != 6 {
		fmt.Print(menu)
		fmt.Scanln(&eleccion)
		scanner := bufio.NewScanner(os.Stdin)
		switch eleccion {
		case 1:
			fmt.Println("Ingresa el nombre:")
			if scanner.Scan() {
				t.User = scanner.Text()
			}
			fmt.Println("Ingresa la Fecha de Creacon:")
			if scanner.Scan() {
				t.Fecha_creacion = scanner.Text()
			}
			fmt.Println("Ingresa la Fecha de Actualizacion:")
			if scanner.Scan() {
				t.Fecha_actualizacion = scanner.Text()
			}
			fmt.Println("Ingresa el Estatus (true, false):")
			if scanner.Scan() {
				t.Estatus, _ = strconv.ParseBool(scanner.Text())
			}
			err := insertar(t)
			if err != nil {
				fmt.Printf("Error insertando: %v", err)
			} else {
				fmt.Println("Insertado correctamente")
			}
		case 2:
			tickets, err := obtenerTickets()
			if err != nil {
				fmt.Printf("Error obteniendo tickets: %v", err)
			} else {
				for _, ticket := range tickets {
					fmt.Println("====================")
					fmt.Printf("Id: %d\n", ticket.Id)
					fmt.Printf("Usuario: %s\n", ticket.User)
					fmt.Printf("Fecha de Creacion: %s\n", ticket.Fecha_creacion)
					fmt.Printf("Fecha de Actualizacion: %s\n", ticket.Fecha_actualizacion)
					fmt.Printf("Estatus: %s\n", ticket.Estatus)
				}
			}

		case 3:
			fmt.Println("Ingresa el Usuario que desea buscar:")
			fmt.Scanln(&t.User)
			err := eliminar(t)
			tickets, err := filtrar(t.User)
			if err != nil {
				fmt.Printf("Error obteniendo tickets: %v", err)
			} else {

				fmt.Println("====================")
				fmt.Printf("Id: %d\n", tickets.Id)
				fmt.Printf("Usuario: %s\n", tickets.User)
				fmt.Printf("Fecha de Creacion: %s\n", tickets.Fecha_creacion)
				fmt.Printf("Fecha de Actualizacion: %s\n", tickets.Fecha_actualizacion)
				fmt.Printf("Estatus: %s\n", tickets.Estatus)
			}

		case 4:
			fmt.Println("Ingresa el id:")
			fmt.Scanln(&t.Id)
			fmt.Println("Ingresa el nuevo nombre:")
			if scanner.Scan() {
				t.User = scanner.Text()
			}
			fmt.Println("Ingresa la nueva fecha de Creacion:")
			if scanner.Scan() {
				t.Fecha_creacion = scanner.Text()
			}
			fmt.Println("Ingresa la nueva fecha de Actualizacion:")
			if scanner.Scan() {
				t.Fecha_actualizacion = scanner.Text()
			}
			fmt.Println("Ingresa el nuevo Estatus:")
			if scanner.Scan() {
				t.Estatus, _ = strconv.ParseBool(scanner.Text())
			}
			err := actualizar(t)
			if err != nil {
				fmt.Printf("Error actualizando: %v", err)
			} else {
				fmt.Println("Actualizado correctamente")
			}
		case 5:
			fmt.Println("Ingresa el ID del ticket que deseas eliminar:")
			fmt.Scanln(&t.Id)
			err := eliminar(t)
			if err != nil {
				fmt.Printf("Error eliminando: %v", err)
			} else {
				fmt.Println("Eliminado correctamente")
			}
		}
	}
}
