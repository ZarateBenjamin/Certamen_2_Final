package main

import (
    "fmt"
    "math/rand"
    "strconv"
)

// dispacher gestiona el despacho de procesos a los núcleos.
func dispacher(NumeroCore int, nombreArchivo string) {
    Termino := 0
    preruta := "./Procesos/" // Define una ruta base para los archivos de procesos
    for Termino < 1 { // Bucle que se ejecuta hasta que Termino sea igual a 1
        CPUTimes[NumeroCore]++ // Incrementa el contador de tiempo de CPU para el núcleo actual
        manejarProcesosPorCrear(NumeroCore, preruta) // Maneja la creación de procesos

        manejarColaListos(NumeroCore) // Maneja la cola de procesos listos
        manejarColaBloqueados(NumeroCore) // Maneja la cola de procesos bloqueados

        // Verifica si las colas y la lista de procesos por crear están vacías
        if len(colaListos) == 0 && len(colaBloqueados) == 0 && len(procesosPorCrear) == 0 {
            Termino = 1 // Si están vacías, establece Termino a 1 para terminar el bucle
        }
    }
    // Prepara y escribe un mensaje indicando que el núcleo terminó su ejecución
    texto := strconv.Itoa(CPUTimes[NumeroCore]) + espacio + "Termino El core " + espacio + strconv.Itoa(NumeroCore) + "\n"
    Escribir(NumeroCore, texto)
}

func manejarProcesosPorCrear(NumeroCore int, preruta string) {
    mutexProcesosPorCrear.Lock() // Bloquea el acceso concurrente a los procesos por crear
    defer mutexProcesosPorCrear.Unlock() // Asegura que se desbloquee al finalizar la función

    // Verifica si hay procesos por crear y si el tiempo de CPU es mayor o igual al tiempo de creación del proceso
    if len(procesosPorCrear) != 0 && (CPUTimes[NumeroCore] >= procesosPorCrear[0].tiempoCreacion) {
        crearProcesoDesdeProcesosPorCrear(NumeroCore, preruta) // Crea un proceso desde la lista de procesos por crear
    }
}

func crearProcesoDesdeProcesosPorCrear(NumeroCore int, preruta string) {
    mutexColaListos.Lock() // Bloquea el acceso concurrente a la cola de procesos listos
    defer mutexColaListos.Unlock() // Asegura que se desbloquee al finalizar la función

    // Crea una nueva instancia de Proceso y la añade a la cola de procesos listos
    proceso := Proceso{
        ID:           ID_Procesos,
        Estado:       "iniciado",
        ContadorPC:   0,
        InfoEstadoES: 0,
        Funcion:      HandleProcessFile,
        nombre:       preruta + procesosPorCrear[0].nombreDelArchivo,
        canalPropio:  make(chan [3]int)}
    procesosPorCrear = procesosPorCrear[1:] // Elimina el proceso de la lista de procesos por crear
    colaListos = append(colaListos, proceso) // Añade el proceso a la cola de procesos listos

    // Prepara y escribe un mensaje indicando la creación del proceso
    texto := strconv.Itoa(CPUTimes[NumeroCore]) + espacio + "PULL" + espacio + "Dispatcher " + espacio + "102" + espacio + strconv.Itoa(NumeroCore) + "\n"
    Escribir(NumeroCore, texto)

    mutexID_Procesos.Lock() // Bloquea el acceso concurrente al ID de procesos
    ID_Procesos++ // Incrementa el ID de procesos
    mutexID_Procesos.Unlock() // Desbloquea el acceso concurrente al ID de procesos
}

func manejarColaListos(NumeroCore int) {
    mutexColaListos.Lock() // Bloquea el acceso concurrente a la cola de procesos listos
    if len(colaListos) > 0 {
        proceso := colaListos[0] // Obtiene el primer proceso de la cola
        colaListos = colaListos[1:] // Elimina el proceso de la cola
        mutexColaListos.Unlock() // Desbloquea el acceso concurrente a la cola de procesos listos

        procesarElementoColaListos(proceso, NumeroCore) // Procesa el elemento obtenido de la cola de procesos listos
    } else {
        mutexColaListos.Unlock() // Desbloquea el acceso concurrente a la cola de procesos listos si está vacía
    }
}

func procesarElementoColaListos(proceso Proceso, NumeroCore int) {
    // Prepara y escribe mensajes indicando la carga y ejecución del proceso
    texto := strconv.Itoa(CPUTimes[NumeroCore]) + espacio + "LOAD" + espacio + proceso.nombre + espacio + "Dispatcher " + espacio + "102" + espacio + strconv.Itoa(NumeroCore) + "\n"
    Escribir(NumeroCore, texto)
    CPUTimes[NumeroCore]++ // Incrementa el tiempo de CPU
    texto = strconv.Itoa(CPUTimes[NumeroCore]) + espacio + "EXEC" + espacio + proceso.nombre + espacio + "Dispatcher " + espacio + "104" + espacio + strconv.Itoa(NumeroCore) + "\n"
    Escribir(NumeroCore, texto)

    // Maneja el estado del proceso
    if proceso.Estado == "iniciado" {
        proceso.Estado = "executando"
        // Inicia la función del proceso en una nueva goroutine
        go proceso.Funcion(proceso.nombre, proceso.ContadorPC, proceso.canalPropio, NumeroCore, (proceso.ID*1000)+5000)
    } else if proceso.Estado == "executando" {
        randomValue := rand.Float64() // Genera un valor aleatorio
        prob := 1.0 / float64(probabilidad) // Calcula la probabilidad
        if randomValue > prob {
            CPUTimes[NumeroCore]++ // Incrementa el tiempo de CPU
            // Prepara y escribe un mensaje indicando la finalización anticipada del proceso
            texto = strconv.Itoa(CPUTimes[NumeroCore]) + espacio + proceso.nombre + " Finalizo de manera anticipada en base al parametro P" + "\n"
            Escribir(NumeroCore, texto)
            return
        }
        proceso.canalPropio <- [3]int{0, 0, NumeroCore} // Envía una señal a través del canal propio del proceso
    }
    senal := <-proceso.canalPropio // Espera una señal de la CPU

    // Maneja la señal recibida
    switch senal[0] {
    case 1:
        // Proceso vuelve a la cola de listos
        CPUTimes[NumeroCore]++ // Incrementa el tiempo de CPU
        mutexColaListos.Lock()
        colaListos = append(colaListos, proceso)
        mutexColaListos.Unlock()
        proceso.Estado = "executando"
        // Prepara y escribe un mensaje indicando que el proceso está en espera
        texto = strconv.Itoa(CPUTimes[NumeroCore]) + espacio + "ST" + espacio + proceso.nombre + espacio + "Dispatcher " + espacio + "100" + espacio + strconv.Itoa(NumeroCore) + "\n"
        Escribir(NumeroCore, texto)
    case 2:
        // Proceso pasa a estado bloqueado
        CPUTimes[NumeroCore]++ // Incrementa el tiempo de CPU
        proceso.InfoEstadoES = senal[1] // Establece la información del estado ES
        proceso.Estado = "bloqueado"
        mutexColaBloqueados.Lock()
        colaBloqueados = append(colaBloqueados, proceso)
        mutexColaBloqueados.Unlock()
        // Prepara y escribe un mensaje indicando que el proceso está bloqueado
        texto = strconv.Itoa(CPUTimes[NumeroCore]) + espacio + "ST" + espacio + proceso.nombre + espacio + "Dispatcher " + espacio + "100" + espacio + strconv.Itoa(NumeroCore) + "no creo esta \n"
        Escribir(NumeroCore, texto)
    case 3:
        // Proceso finaliza
        texto = strconv.Itoa(CPUTimes[NumeroCore]) + espacio + "Finish" + espacio + proceso.nombre + espacio + "Dispatcher " + espacio + "105" + espacio + strconv.Itoa(NumeroCore) + "\n"
        Escribir(NumeroCore, texto)
    default:
        fmt.Println("Esto no deberia pasar nunca") // Mensaje por si ocurre un caso no esperado
    }
}

// manejarColaBloqueados gestiona los procesos en la cola de bloqueados para un núcleo de CPU específico.
func manejarColaBloqueados(NumeroCore int) {
    // Bloquear el acceso a la cola de procesos bloqueados para evitar condiciones de carrera.
    mutexColaBloqueados.Lock()

    // Iterar sobre todos los procesos en la cola de bloqueados.
    for i := 0; i < len(colaBloqueados); i++ {
        // Disminuir el tiempo restante de bloqueo del proceso basado en los ciclos de CPU consumidos.
        colaBloqueados[i].InfoEstadoES -= ciclos

        // Verificar si el proceso ha completado su tiempo de bloqueo.
        if colaBloqueados[i].InfoEstadoES <= 0 {
            // Incrementar el contador de veces que la CPU ha manejado un evento de E/S.
            CPUTimes[NumeroCore]++

            // Registrar el evento de E/S y el movimiento del proceso a la cola de listos.
            texto := strconv.Itoa(CPUTimes[NumeroCore]) + espacio + "** EVENTO E/S - Mover Proceso desde Cola Bloqueado " + colaBloqueados[i].nombre + " a Cola Listo **" + espacio + strconv.Itoa(NumeroCore) + "\n"
            Escribir(NumeroCore, texto)

            // Cambiar el estado del proceso a ejecutando.
            colaBloqueados[i].Estado = "executando"

            // Bloquear el acceso a la cola de listos para agregar el proceso.
            mutexColaListos.Lock()
            colaListos = append(colaListos, colaBloqueados[i])
            mutexColaListos.Unlock()

            // Eliminar el proceso de la cola de bloqueados.
            colaBloqueados = append(colaBloqueados[:i], colaBloqueados[i+1:]...)
            i-- // Ajustar el índice debido a la eliminación de un elemento en la cola.
        }
    }

    // Desbloquear el acceso a la cola de procesos bloqueados.
    mutexColaBloqueados.Unlock()
}

