import { useState, useEffect } from "react";
import { Link, useParams } from "react-router-dom";
import { IoArrowBack } from "react-icons/io5";
import partition from "../../assets/partition.png";
import { ENDPOINT } from "../../App";

function Partition() {
    const [files, setFiles] = useState([]);
    const { driveletter } = useParams();

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch(`${ENDPOINT}/partitions`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: new URLSearchParams({
                        letter: driveletter,
                    }),
                });
                const data = await response.json();
                setFiles(data);
            } catch (error) {
                console.log('Error fetching data:', error);
            }
        };

        fetchData();
    }, [driveletter]);

    return (
        <div className="container-fluid bg-main py-4">
            <div className="container bg-files altoReport">
                <h1 className="display-1 text-light text-center">
                    <Link to={'/files'}>
                        <IoArrowBack />
                    </Link>
                    Selecciona una partición
                </h1>
                <div className="container">
                    {
                        files
                            ? (
                                <ul className='d-flex flex-wrap container-fluid'>
                                    {files.filter(file => file)  // Filtra cualquier valor falso, incluyendo strings vacíos
                                        .map((file, index) => (
                                            <Link key={index} className='nav-link disco p-4' to={`/files/${driveletter}/${file}`}>
                                                <img src={partition} alt="partición" className='img-fluid' />
                                                <h3 className='text-light'>{file}</h3>
                                            </Link>
                                        ))}
                                </ul>

                            )
                            : (<h4>No se encontraron particiones</h4>)
                    }
                </div>
            </div>
        </div>
    );
}

export default Partition;
