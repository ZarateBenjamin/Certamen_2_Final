package main

import (
	"bufio" // Importa el paquete bufio para el manejo de I/O bufferizado
	"fmt"   // Importa el paquete fmt para operaciones de I/O formateadas (por ejemplo, imprimir en consola)
	"log"   // Importa el paquete log para registrar mensajes de log
	"os"    // Importa el paquete os para interactuar con el sistema operativo
	"strconv" // Importa el paquete strconv para la conversión de strings a otros tipos de datos y viceversa
	"strings" // Importa el paquete strings para operaciones con cadenas de texto
)

// Escribir registra los eventos en los archivos de salida correspondientes a cada núcleo.
func Escribir(NumeroCore int, texto string) {
	archivo := archivosSalida[NumeroCore] // Accede al archivo de salida correspondiente al núcleo indicado
	_, err := archivo.WriteString(texto)  // Escribe el texto en el archivo
	if err != nil {
		fmt.Println(err)  // Imprime el error en caso de existir
		return
	}
}

// Creador_Procesos lee un archivo de texto y genera una lista de preprocesos a partir de él.
func Creador_Procesos(nombreArchivo string) []preProceso {
	procesosPorCrear := make([]preProceso, 0) // Inicializa un slice de preProceso
	ArchivoProcesos, err := os.Open(nombreArchivo) // Abre el archivo con el nombre proporcionado
	if err != nil {
		log.Fatal(err) // Si hay un error al abrir el archivo, termina el programa
	}
	defer ArchivoProcesos.Close() // Asegura que el archivo se cierre al finalizar la función

	scanner := bufio.NewScanner(ArchivoProcesos) // Crea un scanner para leer el archivo
	for scanner.Scan() { // Itera sobre cada línea del archivo
		linea := scanner.Text() // Obtiene la línea actual
		if strings.HasPrefix(linea, "#") { // Ignora las líneas que comienzan con "#"
			continue
		}
		texto := strings.Split(linea, " ") // Divide la línea en partes usando el espacio como separador
		numero, _ := strconv.Atoi(texto[0]) // Convierte el primer elemento a entero
		if len(texto) >= 2 { // Verifica si hay más elementos además del número
			for i := 1; i < len(texto); i++ { // Itera sobre los elementos restantes de la línea
				preProcesoVar := preProceso{tiempoCreacion: numero, nombreDelArchivo: texto[i]} // Crea una instancia de preProceso
				procesosPorCrear = append(procesosPorCrear, preProcesoVar) // Agrega el preProceso al slice
			}
		}
	}

	return procesosPorCrear // Retorna el slice de preProceso
}

func HandleProcessFile(nombreArchivo string, cursor int, canal chan [3]int, NumeroCore int, PC int) {
	archivo, err := os.Open(nombreArchivo) // Abre el archivo indicado
	if err != nil {
		log.Fatal(err) // Si hay un error al abrir el archivo, termina el programa
	}
	scanner := bufio.NewScanner(archivo) // Crea un scanner para leer el archivo
	i := 0
	for scanner.Scan() { // Itera sobre cada línea del archivo
		i++
		CPUTimes[NumeroCore]++ // Incrementa el contador de tiempo de CPU para el núcleo indicado
		linea := scanner.Text() // Obtiene la línea actual
		texto := strings.Split(linea, "\t") // Divide la línea en partes usando el tabulador como separador
		if texto[0] == "#" {
			i--
			CPUTimes[NumeroCore]-- // Decrementa el contador si la línea comienza con "#"
			continue
		}
		PCPrograma, _ := strconv.Atoi(texto[0]) // Convierte el primer elemento a entero
		PCTexto := strconv.Itoa(PC + PCPrograma) // Suma el valor convertido a PC y lo convierte a string
		if texto[0] == "Instruccion" {
			texto = strings.Split(texto[0], " ")
			texto := strconv.Itoa(CPUTimes[NumeroCore]) + espacio + "Instruccion" + espacio + nombreArchivo + espacio + PCTexto + "\n"
			Escribir(NumeroCore, texto) // Escribe la instrucción en el archivo de salida
		} else if texto[0] == "Finalizar" {
			respuesta := [3]int{3, 0, 0} // Prepara una respuesta para finalizar
			canal <- respuesta // Envía la respuesta a través del canal
		} else if texto[0] ==  "ES" {
			aux := strings.Split(texto[1], " ")
			texto := strconv.Itoa(CPUTimes[NumeroCore]) + espacio + texto[1] + espacio + nombreArchivo + espacio + PCTexto + "\n"
			Escribir(NumeroCore, texto) // Escribe la operación ES en el archivo de salida
			num, _ := strconv.Atoi(aux[1])
			respuesta := [3]int{2, num, 0} // Prepara una respuesta para ES
			canal <- respuesta // Envía la respuesta a través del canal
			senal := <-canal // Espera una señal a través del canal
			NumeroCore = senal[2] // Actualiza el número de core
			i = senal[0] // Actualiza el índice
		}
		if i == ciclos {
			respuesta := [3]int{1, 0, 0} // Prepara una respuesta para indicar que se alcanzó el número de ciclos
			canal <- respuesta // Envía la respuesta a través del canal
			senal := <-canal // Espera una señal a través del canal
			NumeroCore = senal[2] // Actualiza el número de core
			i = senal[0] // Actualiza el índice
		}

	}
}
