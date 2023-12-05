package main

// Proceso define la estructura para un proceso en el sistema.
type Proceso struct {
    ID           int    // Identificador único del proceso
    Estado       string // Estado actual del proceso (ej: 'iniciado', 'bloqueado', etc.)
    ContadorPC   int    // Contador de Program Counter, indicador de la posición actual en la ejecución del proceso
    InfoEstadoES int    // Información adicional sobre el estado del proceso, posiblemente relacionada con Entrada/Salida (E/S)
    Funcion      func(nombreArchivo string, cursor int, canal chan [3]int, NumeroCore int, PC int) // Función que ejecutará el proceso. La función es pasada como parámetro para permitir diferentes comportamientos.
    nombre       string // Nombre del archivo asociado con el proceso
    canalPropio  chan [3]int // Canal de comunicación propio del proceso, posiblemente para comunicación con CPU o manejo de estados
}

// preProceso define una estructura para los procesos antes de su creación.
type preProceso struct {
    tiempoCreacion   int    // Tiempo en el que el proceso debe ser creado
    nombreDelArchivo string // Nombre del archivo asociado al proceso que se creará
}

