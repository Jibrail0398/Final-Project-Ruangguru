import {ReactComponent as ListIcon } from '../assets/list.svg';
import { useLocation, Link,useNavigate } from 'react-router-dom';
import React, { useState } from "react";
import ChatStart from "../component/ChatStart";
import ChatEnd from "../component/ChatEnd";
import axios from "axios";
import  {ReactComponent as Send} from "../assets/send.svg";


function ChatPage(){

    let navigate = useNavigate();
    const [isDrawerOpen, setDrawerOpen] = useState(false);
    const [file, setFile] = useState(null);
    const [query, setQuery] = useState("");
    const [response, setResponse] = useState("");
    
    const location = useLocation();
    const isActive = (path) => location.pathname === path;

    const handleFileChange = (e) => {
        setFile(e.target.files[0]);
    };

    const handleUpload = async () => {
        const formData = new FormData();
        formData.append("file", file);

        try {
        const res = await axios.post('http://localhost:8080/upload', formData, {
            headers: {
            'Content-Type': 'multipart/form-data',
            },
        });
        console.log(res)
        } catch (error) {
        console.error('Error uploading file:', error);
        }
    };

    

    const Logout = ()=>{
        const drawerToggle = document.getElementById('my-drawer');
        if (drawerToggle) {
            drawerToggle.checked = false;
        } 
        setDrawerOpen(false);
        localStorage.clear()
        navigate("/")
    }
    
    return(
        <>
            <div  className='fixed right-0 top-0 m-4 mb-8  '>
                <button className='bg-red-500  rounded-lg text-white p-4' onClick={Logout} >Logout</button>
            </div>
            <div  >
                <div className="drawer z-[100] w-1/2 ">
                    <input id="my-drawer" type="checkbox" className="drawer-toggle" />
                    <div className="drawer-content ">
                        
                        <label htmlFor="my-drawer">
                            <ListIcon className='  ListIcon cursor-pointer m-2 mb-8' width={50} height={50} />
                        </label>
                    
                    </div>
                    <div className="drawer-side ">
                        <label htmlFor="my-drawer" aria-label="close sidebar" className="drawer-overlay z-100  "></label>
                        <ul className="menu bg-base-200 overflow-y-auto flex-1 text-base-content min-h-full w-80 p-4">
                            <h1 className='text-2xl  font-bold font-roboto mb-5' >History</h1>
                            {/* Sidebar content here */}
                            <Link
                            to="/sidebar-item-1"
                            className={`text-lg mb-2 ${isActive('/Chat') ? 'mb-2 bg-blue-500 text-white rounded-xl p-2' : ''}`}
                            >
                            Sidebar Item 1
                            </Link>

                            <Link
                            to="/sidebar-item-1"
                            className={`text-lg mb-2 ${isActive('/Chat') ? 'mb-2 bg-blue-500 text-white rounded-xl p-2' : ''}`}
                            >
                            Sidebar Item 1
                            </Link>
                           
                            
                            
                        </ul>
                    </div>
                    
                </div>
                
                <div className='  sm:flex sm:justify-center sm:items-center sm:space-x-2'>
                    <input 
                        type="file" 
                        className="mx-auto block sm:w-3/4 max-w-lg p-3 text-sm text-gray-900 border border-gray-300 rounded-lg cursor-pointer bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-700 hover:bg-gray-100 transition-all duration-300 ease-in-out
                        file:mr-4 file:py-2 file:px-4 
                        file:rounded-full file:border-0 
                        file:text-sm file:font-semibold 
                        file:bg-blue-50 file:text-blue-700 
                        hover:file:bg-blue-100 
                        sm:m-0" 
                        onChange={handleFileChange} 
                    />
                    <button 
                        className='block mx-auto p-2 w-3/4 mt-4 sm:w-1/4 sm:mt-0 bg-sky-500 rounded-xl text-white sm:ml-0 sm:px-4' 
                    >
                        Upload
                    </button>

                </div>
            </div>
            <ChatEnd/>
            <ChatStart/>
            <ChatEnd/>
            <ChatStart/>
            <ChatEnd/>
            <ChatStart/>

            <div className="fixed bottom-0 left-0 right-0 w-full p-4 bg-white shadow-lg  flex items-center space-x-2 h-10">
                <input 
                    type="text" 
                    className="input input-bordered  w-full mb-14 p-4" 
                />
                <button >
                    <Send className='Send mb-10 ' width={40} height={40} />
                </button>
            </div>
        </>
    )
}

export default ChatPage;