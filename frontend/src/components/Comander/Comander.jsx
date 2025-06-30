import CodeMirror from '@uiw/react-codemirror'
import { consoleDark } from '@uiw/codemirror-theme-console';
import Swal from 'sweetalert2'
import { BsFillSendFill } from "react-icons/bs";
import { FaUpload } from "react-icons/fa";
import { useEffect, useState } from 'react';
import { ENDPOINT } from '../../App';

function Comander() {
    // const [nombre del elemento, nombre de la funcion que actualiza el elemento] = 
    const [inputValue, setInputValue] = useState('');
    const [terminalValue, setTerminalValue] = useState('');
    const [message, setMessage] = useState('');
    /* -------------------------------------------------------------------------- */
    /*                              MENSAJE DE INICIO                             */
    /* -------------------------------------------------------------------------- */

    /*
        useEffect: es una funcion que actualiza en base a una variable
    */
    useEffect(() => {
        fetch(`${ENDPOINT}/`)
            .then(response => response.text())
            .then(data => setMessage(data));
    }, []);

    /* -------------------------------------------------------------------------- */
    /*                       PERSISTENCIA DE LA INFORMACION                       */
    /* -------------------------------------------------------------------------- */
    useEffect(() => {
        // Cargar el valor guardado cuando el componente se monta
        const savedValue = sessionStorage.getItem('terminalValue');
        if (savedValue) {
            setTerminalValue(savedValue);
        }
    }, []);

    useEffect(() => {
        //LocalStorage: guarda los elementos usando los datos locales (de tu pc)
        //SessionStorage: guarda los elementos hasta que se cierre la pestaña o hasta que se cierre el navegador
        // Guardar el valor en sessionStorage cuando terminalValue cambie
        //sessionStorage.setItem('nombre del item', valor del item)
        sessionStorage.setItem('terminalValue', terminalValue);
    }, [terminalValue]);

    /* -------------------------------------------------------------------------- */
    /*                                 TECLA ENTER                                */
    /* -------------------------------------------------------------------------- */
    const handleKeyDown = (event) => {
        if (event.key === 'Enter') {
            handleSubmit(event);  // Llamar a handleSubmit si la tecla presionada es ENTER
        }
    };

    /* -------------------------------------------------------------------------- */
    /*                       Mandamos a analizar un comando                       */
    /* -------------------------------------------------------------------------- */
    async function handleSubmit() {
        if (!inputValue.trim()) {
            Swal.fire({
                title: 'Error',
                text: 'Cadena vacia',
                icon: 'error'
            });
        } else {
            try {
                const response = await fetch(`${ENDPOINT}/command`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: new URLSearchParams({
                        texto: inputValue,
                    }),
                });
                const data = await response.text();
                setTerminalValue(data);
                setInputValue('')
            } catch (error) {
                Swal.fire({
                    title: 'Error',
                    text: 'No se logro analizar el comando',
                    icon: 'error'
                });
                console.log(error);
            }
        }
    }

    /* -------------------------------------------------------------------------- */
    /*                       Mandamos a analizar un archivo                       */
    /* -------------------------------------------------------------------------- */
    const handleFiles = (event) => {
        const file = event.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = function (event) {

                // Una vez que el archivo es leído, envía su contenido al servidor
                fetch(`${ENDPOINT}/upload`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: new URLSearchParams({
                        fileContent: event.target.result
                    })
                })
                    .then(response => response.text())
                    .then(data => {
                        setTerminalValue(data); // Respuesta del servidor
                    })
                    .catch(error => {
                        Swal.fire({
                            title: 'Error',
                            text: 'No se logro cargar el archivo',
                            icon: 'error'

                        });
                        console.error('Error:', error);
                    });
            };
            reader.readAsText(file);
        }
    };

    return (
        <div className="container-fluid bg-main py-4">
            <div className="container bg-command text-light py-4">
                <h1 className="display-1 text-center">{message}</h1>
                <div className="container">
                    <CodeMirror
                        width='100%'
                        height='60vh'
                        readOnly='true'
                        value={terminalValue}
                        theme={consoleDark}
                    />
                    <br />
                    <div className='d-flex'>
                        <input
                            type="text"
                            className="form-control anchoInputComando"
                            placeholder="Ingresar comando"
                            value={inputValue}
                            onChange={(e) => setInputValue(e.target.value)}
                            onKeyDown={handleKeyDown}
                        />
                        <button type="button" onClick={() => handleSubmit()} className="btn btn-primary anchoBtnComando"><BsFillSendFill /></button>
                        {/* Codigo para hacer que la busqueda de archivos en un boton */}
                        <label htmlFor="fileInput" className="btn btn-danger anchoBtnComando">
                            <FaUpload />
                        </label>
                        <input
                            id="fileInput"
                            type="file"
                            onChange={handleFiles}
                            className="d-none" // Hacer el input de archivo invisible
                        />
                    </div>
                    <br />
                    <hr />
                    <h2 className='display-2 text-center'>Comandos</h2>
                    <div className='fs-6'>
                        <ul>
                            <strong className='fs-3'>Administracion de discos</strong>
                            <li><strong>MKDISK: </strong>Crea un nuevo disco.</li>
                            <li><strong>RMDISK: </strong>Elimina un disco</li>
                            <li><strong>FDISK: </strong>Crea una particion (P/L/E)</li>
                            <li><strong>MOUNT: </strong>Monta una particion</li>
                            <li><strong>UNMOUNT: </strong>Desmonta una particion</li>
                            <li><strong>MKFS: </strong>Aplica el formato EXT2/EXT3</li>
                            <br />
                            <strong className='fs-3'>Administracion de usuarios y grupos</strong>
                            <li><strong>LOGIN: </strong>Inicio de sesion</li>
                            <li><strong>LOGOUT: </strong>Cerrar sesion</li>
                            <li><strong>MKGRP: </strong>Crea un grupo</li>
                            <li><strong>RMGRP: </strong>Elimina un grupo</li>
                            <li><strong>MKUSR: </strong>Crea un usuario</li>
                            <li><strong>RMUSR: </strong>Elimina un usuario</li>
                            <li><strong>CHGRP: </strong>Cambia de grupo</li>
                            <br />
                            <strong className='fs-3'>Administracion de carpetas y permisos.</strong>
                            <li><strong>MKFILE: </strong>Creara un archivo en la ruta especifica.</li>
                            <li><strong>CAT: </strong>Muestra el contenido de un archivo.</li>
                            <li><strong>REMOVE: </strong>Eliminara un archivo o carpeta.</li>
                            <li><strong>EDIT: </strong>Edita el contenido de un archivo.</li>
                            <li><strong>RENAME: </strong>Cambia el nombre de un archivo o carpeta.</li>
                            <li><strong>MKDIR: </strong>Creara una carpeta en la ruta especificada.</li>
                            <li><strong>COPY: </strong>Copia un archivo o carpeta.</li>
                            <li><strong>MOVE: </strong>Cambia la ubicacion de un archivo o carpeta.</li>
                            <li><strong>FIND: </strong>Realiza la busqueda de un archivo o carpeta.</li>
                            <li><strong>CHOWN: </strong>Cambiara de propietario el archivo o carpeta.</li>
                            <li><strong>CHMOD: </strong>Cambiara los permisos de un usuario en la particion.</li>
                            <li><strong>PAUSE: </strong>Detendra los procesos.</li>
                            <br />
                            <strong className='fs-3'>Perdida y recuperacion del sistema de archivos EXT3</strong>
                            <li><strong>RECOVERY: </strong>Recobrara los archivos y carpetas en caso de perdida.</li>
                            <li><strong>LOSS: </strong>Formatea los registros.</li>
                            <br />
                            <strong className='fs-3'>Comentarios</strong>
                            <li><strong>#: </strong>los comentarios no se mostraran en consola</li>
                            <br />
                            <strong className='fs-3'>Reportes</strong>
                            <li><strong>REP: </strong>Generara el reporte indicado.</li>


                        </ul>
                    </div>
                </div>
            </div>

        </div>
    )
}

export default Comander;