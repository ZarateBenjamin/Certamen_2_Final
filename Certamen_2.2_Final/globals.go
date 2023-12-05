package main

import (
	"os"    // Importación del paquete os para manejo de archivos
	"sync"  // Importación del paquete sync para utilizar mutex y sincronizar goroutines
)

// Variables globales
var (
	CPUTimes           []int // Slice para almacenar tiempos de CPU para cada núcleo
	probabilidad       int = 0 // Variable para definir la probabilidad de un evento específico (no especificado aquí)
	ciclos             int = 0 // Contador de ciclos de CPU o similar
	archivosSalida     []*os.File // Slice para manejar múltiples archivos de salida
	colaListos         = make([]Proceso, 0)     // Cola de procesos listos para ejecución
	colaBloqueados     = make([]Proceso, 0)     // Cola de procesos bloqueados
	procesosPorCrear   = make([]preProceso, 0)  // Lista de procesos pendientes de creación
	ID_Procesos        = 1                      // Contador de ID para los procesos, inicia en 1
	mutexColaListos    sync.Mutex               // Mutex para sincronizar el acceso a colaListos
	mutexColaBloqueados sync.Mutex              // Mutex para sincronizar el acceso a colaBloqueados
	mutexProcesosPorCrear sync.Mutex            // Mutex para sincronizar el acceso a procesosPorCrear
	mutexID_Procesos   sync.Mutex               // Mutex para sincronizar la asignación de ID_Procesos
)

const espacio string = "    " // Constante para definir un espacio, usado para la escritura en archivos
