import { ReactComponent as MyIcon } from '../assets/visionary.svg';
import { Link } from 'react-router-dom';
import { useState,useRef } from 'react';
import axios from 'axios';
import Loader from "../component/Loader"
import Modal from "../component/Modal"

function RegisterPage(){

    const[email,setEmail] = useState("");
    const[password,setPassword] = useState("");
    const[username,setUsername] = useState("");
    const[isLoading,setIsloading] = useState(false);
    

    const modalRef = useRef();
    const [modalContent, setModalContent] = useState({
        header: "",
        message: "",
    });
    

    const handleRegister = async (e)=>{
        e.preventDefault();


        const requestBody = {
            username:username,
            email:email,
            password:password,
        }

        try{
            setIsloading(true);
            const response = await axios.post("http://localhost:8080/register",requestBody,{
                headers:{
                    "Content-Type":"application/json",
                }
            });
            console.log(response)
            setIsloading(false)
            setModalContent({
                header: "Success",
                message: `Login Success`,
            });
            if (modalRef.current) {
                modalRef.current.openModal();
            }

        }catch(error){
            console.log("Error",error)
            setIsloading(false)
            setModalContent({
                header:"Login Failed",
                message:error.response.data.message,
            })
            if (modalRef.current) {
                modalRef.current.openModal();
            }
            
        }
    }

    return(
        <>
            
            <Modal ref={modalRef} header={modalContent.header} message={modalContent.message} />
            <div className="flex justify-center items-center h-screen " >

                <div className="card bg-base-100 w-full h-full md:w-3/4 md:h-1/2 lg:h-full shadow-xl">
                    <div className="card-body  lg:flex-row ">

                        <div className=' lg:w-1/2' >

                            <h2 className="card-title text-3xl ">Register!</h2>
                            <label className="form-control w-full max-w-xs">
                                <div className="label">
                                    <span className="label-text text-lg">Username</span>
                                </div>
                                <input 
                                onChange={(e)=>setUsername(e.target.value)}
                                type="text" placeholder="Type here" className="input input-bordered w-full max-w-xs" />
                            </label>
                            <label className="form-control w-full max-w-xs">
                                <div className="label">
                                    <span className="label-text text-lg">Email</span>
                                </div>
                                <input 
                                type="email" 
                                onChange={(e)=>setEmail(e.target.value)}
                                
                                placeholder="Type here" className="input input-bordered w-full max-w-xs" />
                            </label>
                            <label className="form-control w-full max-w-xs">
                                <div className="label">
                                    <span className="label-text text-lg">Password</span>
                                </div>
                                <input 
                                onChange={ (e)=> setPassword(e.target.value)}
                                type="password" placeholder="Type here" className="input input-bordered w-full max-w-xs" />
                            </label>
                        </div>

                        <div  className='mx-auto self-center hidden  md:hidden lg:block '>
                            <MyIcon className="my-icon" style={{ width: '300px', height: '300px' }} />
                        </div>

                    </div>
                        <Link to="/" >
                            <p className='text-center' >Already have an Account? <span className='text-primary' >Log in Here</span></p>
                        </Link>
                        
                    <div className="card-actions justify-end">

                        <button onClick={handleRegister} className="btn btn-success w-3/4 mx-auto m-2 text-white ">
                        
                        {isLoading ? <Loader/>:"Register"}
                        </button>
                    </div>
                </div>
            </div>
        </>
    )
}

export default RegisterPage;