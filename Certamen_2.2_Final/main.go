package main
//go run . 3 100 0.1 orden_ejecucion.txt salida.txt
// Comentario que indica cómo ejecutar el programa desde la línea de comandos.
// Los argumentos son: número de núcleos, ciclos máximos, probabilidad, archivo de orden de ejecución y nombre de archivo de salida.

import (
    "fmt" // Importa el paquete fmt para operaciones de entrada/salida formatas, como imprimir en la consola.
    "math/rand" // Importa el paquete rand para generar números aleatorios.
    "os" // Importa el paquete os para interactuar con el sistema operativo, como manejo de archivos y argumentos de la línea de comandos.
    "strconv" // Importa el paquete strconv para convertir strings a otros tipos de datos.
    "time" // Importa el paquete time para manipular fechas y horas.
)

func main() {
    rand.Seed(time.Now().UnixNano()) // Inicializa el generador de números aleatorios con una semilla basada en la hora actual.

    args := os.Args // Almacena los argumentos de la línea de comandos en una variable llamada args.
    if len(args) < 6 { // Verifica si se han proporcionado menos de 6 argumentos.
        fmt.Println("Uso: Orden_Ejecucion_Prog n o P archivo_orden_creacion_procesos nombre_archivo_salida") // Imprime un mensaje de error si no hay suficientes argumentos.
        return // Termina la ejecución del programa si no hay suficientes argumentos.
    }

    NCores, _ := strconv.Atoi(args[1]) // Convierte el segundo argumento a un número entero que representa el número de núcleos.
    ciclos, _ = strconv.Atoi(args[2]) // Convierte el tercer argumento a un número entero que representa los ciclos máximos. Asumiendo que 'ciclos' está declarado en otro lugar.
    probabilidad, _ = strconv.Atoi(args[3]) // Convierte el cuarto argumento a un número entero que representa la probabilidad. Asumiendo que 'probabilidad' está declarada en otro lugar.
    archivoOrden := args[4] // Almacena el quinto argumento como el nombre del archivo de orden de ejecución.
    nombreArchivoSalida := args[5] // Almacena el sexto argumento como el nombre del archivo de salida.

    //se crean los procesos
    procesosPorCrear = Creador_Procesos(archivoOrden) // Llama a una función para crear procesos basados en el archivo de orden. Asumiendo que 'procesosPorCrear' está declarado en otro lugar.
    texto := "# Tiempo de CPU Tipo Instrucción Proceso/Despachador Valor CP\n" // Define un encabezado para los archivos de salida.
    for i := 0; i < NCores; i++ { // Itera sobre el número de núcleos.
        archivoSalida, err := os.Create("Core_" + strconv.Itoa(i) + "_" + nombreArchivoSalida) // Crea un archivo de salida para cada núcleo.
        if err != nil { // Verifica si hubo un error al crear el archivo.
            fmt.Println(err) // Imprime el error en la consola.
            return // Termina la ejecución si hay un error.
        }
        defer archivoSalida.Close() // Asegura que el archivo se cierre cuando la función main termine.
        archivoSalida.WriteString(texto) // Escribe el encabezado en el archivo de salida.
        CPUTimes = append(CPUTimes, 0) // Añade un contador de tiempo de CPU para cada núcleo. Asumiendo que 'CPUTimes' está declarado en otro lugar.
        archivosSalida = append(archivosSalida, archivoSalida) // Añade el archivo de salida a una lista. Asumiendo que 'archivosSalida' está declarada en otro lugar.
        go dispacher(i, nombreArchivoSalida) // Inicia un despachador en una nueva goroutine para cada núcleo.
    }
    // La siguiente línea parece comentada y por lo tanto no se ejecuta:
    // go comerArchivo("./Procesos/proceso_1", 10, canal, archivoSalida) // Ejecuta una función en una goroutine que podría estar procesando un archivo.
    
    time.Sleep(1 * time.Second) // Pausa la ejecución del programa principal por 1 segundo.
    // Las siguientes líneas parecen comentadas y por lo tanto no se ejecutan:
    // valor := <-canal // Recibe un valor de un canal. Asumiendo que 'canal' está declarado en otro lugar.
    // fmt.Printf("%d\n", valor) // Imprime el valor recibido del canal.
    fmt.Print("Fin Simulacion\n") // Imprime un mensaje indicando el fin de la simulación.
}

// PrepararArchivoSalida prepara los archivos de salida para cada núcleo.
func PrepararArchivoSalida(NCores int, nombreArchivoSalida string) {
    const header = "# Tiempo de CPU Tipo Instrucción Proceso/Despachador Valor CP\n" // Define un encabezado constante para los archivos de salida.
    for i := 0; i < NCores; i++ { // Itera sobre el número de núcleos.
        archivoSalida, err := os.Create(fmt.Sprintf("Core_%d_%s", i, nombreArchivoSalida)) // Crea un archivo de salida para cada núcleo con un nombre formateado.
        if err != nil { // Verifica si hubo un error al crear el archivo.
            fmt.Println(err) // Imprime el error en la consola.
            return // Termina la ejecución si hay un error.
        }
        defer archivoSalida.Close() // Asegura que el archivo se cierre cuando la función termine.
        archivoSalida.WriteString(header) // Escribe el encabezado en el archivo de salida.
        CPUTimes = append(CPUTimes, 0) // Añade un contador de tiempo de CPU para cada núcleo.
        archivosSalida = append(archivosSalida, archivoSalida) // Añade el archivo de salida a una lista.
    }
}
