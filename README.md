# api_go

### Instalación
* Clona este repositorio
* Instala las dependencias con:
  ```
    go test
    
  ``` 

* Ejecute este SCRIPT de MySQL:
  ```
  CREATE DATABASE tickets;
  USE tickets;

  CREATE TABLE tickets_guardados (
    id_ticket INT AUTO_INCREMENT, 
    usuario varchar(20),
    fecha_creacion varchar(20),
    fecha_actualizacion varchar(20),
    estatus tinyint(1),
    PRIMARY KEY (id_ticket)
    );
  ```
* Configura las credenciales de acceso a la base de datos en la función ```obtenerBaseDeDatos``` almacenada en el archivo main.go:
  ```
    func obtenerBaseDeDatos() (db *sql.DB, e error) {
      //Cambiar datos para ejecución
      usuario := "root"
      pass := ""
      host := "tcp(localhost)"
      nombreBaseDeDatos := "tickets"
      // Debe tener la forma usuario:contraseña@host/nombreBaseDeDatos
      db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", usuario, pass, host, nombreBaseDeDatos))
      if err != nil {
        return nil, err
      }
      return db, nil
    }
  ```
* Ejecuta el archivo main.go:
  ```
  go run main.go  
  ```
