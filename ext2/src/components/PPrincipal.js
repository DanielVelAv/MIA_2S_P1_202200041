import axios from 'axios';
import React,{useState} from 'react'


export const PPrincipal = () => {

  const [fileContent, setFileContent] = useState("");

  const handleUpload = (e) => {
    e.preventDefault(); // Prevenir el comportamiento por defecto del botón
    const file = document.getElementById("arch").files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = function(event) {
        setFileContent(event.target.result);
      };
      reader.readAsText(file);
    } else {
      alert("Por favor selecciona un archivo primero.");
    }
  };

  const handleTextArea = (e) => {
    setFileContent(e.target.value);
  };

  const handleSubmit = (e) => {
  
    e.preventDefault();
    axios.post('http://localhost:8080/api/process', { text: fileContent })
      .then(response => {
        console.log(response.data);
      })

      .catch(error => {
        console.error("Ocurrio un error",error);
      });

  };


  return (
    <div>
        <form id="sectEnt" onSubmit={handleSubmit}>
          <h1 id="txtEnt">Entrada:</h1>
          <textarea name="entrada" id="entrada" cols="100" rows="10" value={fileContent} onChange={handleTextArea}></textarea> <br></br><br></br><br></br>
          <input type="file" id="arch" />
          <button id="btnEjecutar" onClick={handleUpload}>Subir texto</button>
          <button id="btnEjecutarCom" type='submit'> Ejecutar </button>

        </form>

        <section id="sectSal">
        <h1 id="txtSal">Salida:</h1>
          <textarea name="salida" id="salida" cols="100" rows="10"></textarea> <br></br><br></br><br></br>
        </section>
    </div>
  )
}