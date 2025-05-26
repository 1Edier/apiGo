//utils/exel.go
package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
	"server/models"
)


const (
	ExcelFilePath = "./bdexcel.xlsx" 
	SheetName     = "Hoja1"          
)


func InitExcel() error {
	
	if _, err := os.Stat(ExcelFilePath); os.IsNotExist(err) {
		
		
		fmt.Println("Archivo Excel noo encontrado")
	} 
	
	return nil
}

func GetClientes() ([]models.Cliente, error) {
	f, err := excelize.OpenFile(ExcelFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	
	// Obtener todas las filas
	rows, err := f.GetRows(SheetName)
	if err != nil {
		return nil, err
	}
	
	clientes := []models.Cliente{}
	
	
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) >= 4 { 
			cliente := models.Cliente{
				Clave:            row[0],
				NombreContacto:   row[1],
				Correo:           row[2],
				TelefonoContacto: row[3],
			}
			clientes = append(clientes, cliente)
		}
	}
	
	return clientes, nil
}


func BuscarClientes(termino string) ([]models.Cliente, error) {
	clientes, err := GetClientes()
	if err != nil {
		return nil, err
	}
	
	terminoLower := strings.ToLower(termino)
	resultados := []models.Cliente{}
	
	for _, cliente := range clientes {
		
		if strings.Contains(strings.ToLower(cliente.Clave), terminoLower) ||
		   strings.Contains(strings.ToLower(cliente.NombreContacto), terminoLower) ||
		   strings.Contains(strings.ToLower(cliente.Correo), terminoLower) ||
		   strings.Contains(strings.ToLower(cliente.TelefonoContacto), terminoLower) {
			resultados = append(resultados, cliente)
		}
	}
	
	return resultados, nil
}


func GetClientePorClave(clave string) (*models.Cliente, error) {
	clientes, err := GetClientes()
	if err != nil {
		return nil, err
	}
	
	for _, cliente := range clientes {
		if cliente.Clave == clave {
			return &cliente, nil
		}
	}
	
	return nil, nil 
}


func CrearCliente(cliente models.Cliente) error {
	
	existente, err := GetClientePorClave(cliente.Clave)
	if err != nil {
		return err
	}
	
	if existente != nil {
		return fmt.Errorf("ya existe un cliente con la clave %s", cliente.Clave)
	}
	
	
	f, err := excelize.OpenFile(ExcelFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	
	// Obtener la última fila
	rows, err := f.GetRows(SheetName)
	if err != nil {
		return err
	}
	
	nextRow := len(rows) + 1
	
	// Agregar el nuevo cliente
	f.SetCellValue(SheetName, fmt.Sprintf("A%d", nextRow), cliente.Clave)
	f.SetCellValue(SheetName, fmt.Sprintf("B%d", nextRow), cliente.NombreContacto)
	f.SetCellValue(SheetName, fmt.Sprintf("C%d", nextRow), cliente.Correo)
	f.SetCellValue(SheetName, fmt.Sprintf("D%d", nextRow), cliente.TelefonoContacto)
	
	// Guardar los cambios
	if err := f.SaveAs(ExcelFilePath); err != nil {
		return err
	}
	
	return nil
}

// ActualizarCliente actualiza un cliente existente
func ActualizarCliente(clave string, nuevoCliente models.Cliente) error {
	// Verificar que el cliente existe
	existente, err := GetClientePorClave(clave)
	if err != nil {
		return err
	}
	
	if existente == nil {
		return fmt.Errorf("no existe un cliente con la clave %s", clave)
	}
	
	// Si la clave nueva es diferente, verificar que no exista
	if clave != nuevoCliente.Clave {
		otroExistente, err := GetClientePorClave(nuevoCliente.Clave)
		if err != nil {
			return err
		}
		
		if otroExistente != nil {
			return fmt.Errorf("ya existe otro cliente con la clave %s", nuevoCliente.Clave)
		}
	}
	
	// Abrir el archivo Excel
	f, err := excelize.OpenFile(ExcelFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	
	// Buscar la fila del cliente
	rows, err := f.GetRows(SheetName)
	if err != nil {
		return err
	}
	
	filaEncontrada := 0
	for i := 1; i < len(rows); i++ {
		if len(rows[i]) > 0 && rows[i][0] == clave {
			filaEncontrada = i + 1 // +1 porque los índices en Excel empiezan en 1
			break
		}
	}
	
	if filaEncontrada == 0 {
		return fmt.Errorf("no se encontró el cliente en el archivo Excel")
	}
	
	// Actualizar los datos
	f.SetCellValue(SheetName, fmt.Sprintf("A%d", filaEncontrada), nuevoCliente.Clave)
	f.SetCellValue(SheetName, fmt.Sprintf("B%d", filaEncontrada), nuevoCliente.NombreContacto)
	f.SetCellValue(SheetName, fmt.Sprintf("C%d", filaEncontrada), nuevoCliente.Correo)
	f.SetCellValue(SheetName, fmt.Sprintf("D%d", filaEncontrada), nuevoCliente.TelefonoContacto)
	
	// Guardar los cambios
	if err := f.SaveAs(ExcelFilePath); err != nil {
		return err
	}
	
	return nil
}

// EliminarCliente elimina un cliente
func EliminarCliente(clave string) error {
	// Verificar que el cliente existe
	existente, err := GetClientePorClave(clave)
	if err != nil {
		return err
	}
	
	if existente == nil {
		return fmt.Errorf("no existe un cliente con la clave %s", clave)
	}
	
	// Obtener todos los clientes excepto el que se va a eliminar
	clientes, err := GetClientes()
	if err != nil {
		return err
	}
	
	clientesFiltrados := []models.Cliente{}
	for _, c := range clientes {
		if c.Clave != clave {
			clientesFiltrados = append(clientesFiltrados, c)
		}
	}
	
	// Crear un nuevo archivo Excel
	f := excelize.NewFile()
	sheetIndex, err := f.NewSheet(SheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(sheetIndex)
	
	// Establecer encabezados
	headers := []string{"Clave", "NombreContacto", "Correo", "TelefonoContacto"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(SheetName, cell, header)
	}
	
	// Agregar los clientes filtrados
	for i, cliente := range clientesFiltrados {
		row := i + 2 // +2 porque la fila 1 son los encabezados
		f.SetCellValue(SheetName, fmt.Sprintf("A%d", row), cliente.Clave)
		f.SetCellValue(SheetName, fmt.Sprintf("B%d", row), cliente.NombreContacto)
		f.SetCellValue(SheetName, fmt.Sprintf("C%d", row), cliente.Correo)
		f.SetCellValue(SheetName, fmt.Sprintf("D%d", row), cliente.TelefonoContacto)
	}
	
	// Guardar el archivo
	if err := f.SaveAs(ExcelFilePath); err != nil {
		return err
	}
	
	return nil
}