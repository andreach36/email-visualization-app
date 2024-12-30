# **Email Database Visualization App**

## **Descripción del proyecto**
Esta aplicación permite visualizar y buscar información de una base de datos de correos electrónicos a través de una interfaz intuitiva basada en una tabla, barra de búsqueda y panel de visualización.

La data empleada corresponde al dataset de correos electrónicos de la Corporación Enron, compuesto por aproximadamente 500,000 emails generados por los empleados de la empresa. Esta información fue obtenida por la Comisión Federal Reguladora de la Energía durante su investigación sobre la quiebra de Enron.

## **Caraterísticas principales**
- **Visualización del Dataset**: Presentación de datos en una tabla con paginación de 10 filas por página.
- **Búsqueda por palabra clave**: Permite buscar correos electrónicos ingresando palabras clave en una barra de búsqueda.
- **Panel de información de email**: Muestra detalles como el asunto y el contenido del correo seleccionado.

## **Tecnologías usadas**

### **Frontend**:
- **Vue.js**: Framework para el desarrollo de interfaces de usuario.  
- **TypeScript**: Lenguaje tipado para mayor robustez en el código.  
- **Tailwind CSS**: Framework para diseño de estilos modernos y responsivos.

### **Backend**:  
- **Golang**: Para la indexación de datos y la gestión de la API REST que sirve los datos al frontend.

### **Base de Datos**:  
- **ZincSearch**: Utilizado para almacenar, indexar y buscar datos de forma eficiente.

## **Instrucciones de uso**

1. **Visualización inicial**:
 Al iniciar la aplicación, se muestra una tabla con la lista paginada de correos electrónicos. Cada página contiene 10 registros.

2. **Visualización detallada**:
Haz clic en cualquier fila de la tabla para ver más detalles sobre el asunto y el contenido del correo en el panel lateral derecho.

3. **Búsqueda**:
Ingresa una palabra clave en la barra de búsqueda para filtrar los correos electrónicos relacionados.

## **Capturas de pantalla**
![email data1](https://github.com/user-attachments/assets/870b3133-27de-413d-af48-af1e6ff48f3a)

![email data2](https://github.com/user-attachments/assets/6f1b45c2-0000-4771-8515-0e45eace4374)

![email data3](https://github.com/user-attachments/assets/5dc846f0-3456-4452-9b1d-0cfb58934e68)

![email data4](https://github.com/user-attachments/assets/304701be-7c0f-4fe4-843f-15b9f2b01eaa)





