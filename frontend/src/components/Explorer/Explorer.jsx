import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import Swal from 'sweetalert2';
import { useSession } from '../../session/useSession';
import { ENDPOINT } from '../../App';
import txt from "../../assets/txt-file.png"
import carpeta from "../../assets/folder.png"
function Explorer() {
    const [searchTerm, setSearchTerm] = useState('/');
    const [files, setFiles] = useState([]);
    const navigate = useNavigate()
    const { isAuthenticated } = useSession();
    const { driveletter } = useParams()
    const { partition } = useParams()


    // Mapeo de extensiones a imágenes
    const fileIcons = {
        txt: txt
    };
    // Funcion para obtener la extension
    function getFileExtension(filename) {
        if (typeof filename !== 'string') return '';

        const lastDotIndex = filename.lastIndexOf(".");
        if (lastDotIndex === -1 || lastDotIndex === 0 || lastDotIndex === filename.length - 1) return '';

        const extension = filename.substring(lastDotIndex + 1).toLowerCase();

        // Comprueba si la extensión es 'txt'
        if (extension.endsWith("txt")) {
            return "txt"
        }

        return extension;
    }


    useEffect(() => {
        if (!isAuthenticated) {
            // Mostrar el mensaje de alerta
            Swal.fire({
                title: 'Sesión no válida',
                text: 'Serás redirigido al login',
                icon: 'warning',
                confirmButtonText: 'OK'
            }).then((result) => {
                // Cuando el usuario hace clic en "OK", navegar a la página de inicio
                if (result.isConfirmed) {
                    navigate(`/files/${driveletter}/${partition}`);
                }
            });
        }
    }, [isAuthenticated, navigate, driveletter, partition]);


    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch(`${ENDPOINT}/docs`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        disk: driveletter,
                        partition: partition,
                        ruta: searchTerm
                    }),
                });
                const data = await response.json();
                setFiles(data);
            } catch (error) {
                console.log('Error fetching data:', error);
            }
        };
        
        fetchData()
    }, [driveletter, partition, searchTerm]);
    


    function IncRuta(name) {
        // Comprobar si el searchTerm actual termina con '/'
        if (searchTerm.endsWith('/')) {
            setSearchTerm(searchTerm + name);
        } else {
            setSearchTerm(searchTerm + '/' + name);
        }
    }



    return (
        <div className="container-fluid bg-main py-4">
            <div className="container bg-files altoExplorer">
                <br />
                <input
                    type="text"
                    className="form-control"
                    placeholder="Ingresa la ruta aqui"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                />
                {
                    files
                        ? (
                            <ul className='d-flex flex-wrap container-fluid'>
                                {files.filter(file => file).map((file, index) => {
                                    const ext = getFileExtension(file).toLowerCase(); // Obtener la extensión del archivo
                                    const icon = fileIcons[ext] || carpeta; // Obtener el ícono basado en la extensión o png por defecto
                                    return (
                                        <button onClick={() => IncRuta(file)} key={index} className='nav-link report p-4' to={'/'}>
                                            <img src={icon} alt="archivo" className='img-fluid' />
                                            <h6 className='text-light text-center'>{file}</h6>
                                        </button>
                                    );
                                })}
                            </ul>

                        )
                        : (<h3 className="text-light">No se encontraron archivos</h3>)
                }

            </div>
        </div>
    );
}

export default Explorer;
