import { useState, useEffect } from 'react';
import { ENDPOINT } from '../../App';
import { Link } from 'react-router-dom';
import disco from "../../assets/disk.png"
import Swal from 'sweetalert2';

function File() {
    const [files, setFiles] = useState([]);

    useEffect(() => {
        // Llamada a la API para obtener los nombres de los archivos
        fetch(`${ENDPOINT}/disks`)
            .then(response => {
                if (!response.ok) {
                    Swal.fire({
                        title: 'Error',
                        text: 'Problema de red',
                        icon: 'error'
                    })
                }
                return response.json();
            })
            .then(data => setFiles(data))
            .catch(error => {
                Swal.fire({
                    title: 'Error',
                    text: error,
                    icon: 'error'
                })
            });
    }, []);


    return (
        <div className="container-fluid bg-main py-4">
            <div className="container bg-files altoReport">
                <h1 className="display-1 text-light text-center">Selecciona un disco</h1>
                <div className="container">
                    {
                        files
                            ? (
                                <ul className='d-flex flex-wrap container-fluid'>
                                    {files.map((file, index) => (
                                        // Quitamos el .dsk
                                        file = file.replace(/\.dsk$/, ''),
                                        <Link key={index} className='nav-link disco p-4'  to={`/files/${file}`}>
                                            <img src={disco} alt="disco" className='img-fluid' />
                                            <h3 className='text-light'>{file}</h3>
                                        </Link>
                                    ))}
                                </ul>

                            )
                            : (<h3 className="text-light">No se han creado discos</h3>)
                    }
                </div>
            </div>
        </div>
    )
}

export default File;