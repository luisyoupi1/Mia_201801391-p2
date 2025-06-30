import { useEffect, useState } from 'react';
import { Graphviz } from 'graphviz-react';
import { ENDPOINT } from '../../App';
// import pdf from "../../assets/pdf-file.png"
// import png from "../../assets/png-file.png"
// import svg from "../../assets/svg-file.png"
// import txt from "../../assets/txt-file.png"

// Mapeo de extensiones a imágenes
// const fileIcons = {
//     pdf: pdf,
//     png: png,
//     svg: svg,
//     txt: txt
// };
// Funcion para obtener la extension
// function getFileExtension(filename) {
//     return filename.slice(((filename.lastIndexOf(".") - 1) >>> 0) + 2);
// }

function Report() {

    const [files, setFiles] = useState([]);

    useEffect(() => {
        // Llamada a la API para obtener los nombres de los archivos
        fetch(`${ENDPOINT}/reports`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => setFiles(data))
            .catch(error => {
                console.error('There was a problem with the fetch operation:', error);
            });
    }, []);



    return (
        <div className="container-fluid bg-main py-4">
            <div className="container bg-report altoReport py-4">
                <h1 className="display-1 text-light text-center">Reportes</h1>
                {
                    files
                        ? (
                            <ul className='d-flex flex-wrap container-fluid flex-column'>
                                {files.map((file, index) => {
                                    return (
                                        <li key={index}>
                                            <hr className='border-warning border border-2'/>
                                            <div className='container-fluid'>
                                                <h3 className='text-light text-center display-3'>{file.nombre}</h3>
                                                <Graphviz dot={file.dot} options={{ width:'100%', fit: true, zoom: true }} />
                                            </div>
                                        </li>
                                    );
                                })}
                            </ul>

                        )
                        : (<h3 className="text-light">No se han generado reportes</h3>)
                }
            </div>
        </div>
    )
}

export default Report