import {ReactComponent as ListIcon } from '../assets/list.svg';
import {ReactComponent as DeleteIcon} from '../assets/trash-fill.svg'
import { useLocation, Link,useNavigate } from 'react-router-dom';
import React, { useState,useEffect,useRef } from "react";
import ChatStart from "../component/ChatStart";
import ChatEnd from "../component/ChatEnd";
import axios from "axios";
import  {ReactComponent as Send} from "../assets/send.svg";
import Modal from "../component/Modal"
import Loader from "../component/Loader"

function ChatPage(){

    let navigate = useNavigate();
    const modalRef = useRef();
    const [file, setFile] = useState(null);
    const [reports,setReport] = useState(null);
    const [loading,setLoading] = useState(false);
    const [question,setQuestion] = useState("");

    // Referensi untuk container chat
    const chatContainerRef = useRef(null);

    const [aiResponse,setAIResponse] = useState({
        question:[],
        answer:[],
    });
    const [isLoadingChat,setIsLoadingChat] = useState(false)
    
    
    const [hasUploadFile,setHasUploafFile] = useState(false);
    const [uploadFileData,setUploadFileData] = useState(null);
    const [modalContent, setModalContent] = useState({
        header: "",
        message: "",
    });
    const id_user = localStorage.getItem("id_user")
    const token = localStorage.getItem("token")
    
    const reset = ()=>{
        setFile(null)
        setQuestion("")
        setAIResponse({
            question:[],
            answer:[],
        });
        setHasUploafFile(false)
        setUploadFileData(null)

    }

    const handleFileChange = (e) => {
        setFile(e.target.files[0]);
    };

    useEffect(()=>{
        getReport()

        return()=>{
            
            localStorage.removeItem("id_report")
            localStorage.removeItem("date")
        }
        
    },[])

    // Fungsi untuk menggulir ke bawah
    const scrollToBottom = () => {
        if (chatContainerRef.current) {
            chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
        }
    };

    // Scroll otomatis setiap kali ada pesan baru
    useEffect(() => {
        scrollToBottom();
    }, [aiResponse]);

    const handleUpload = async () => {

        setLoading(true)
        const formData = new FormData();

        formData.append("file", file);
        formData.append("id_user", id_user);


        try {
        const res = await axios.post('http://localhost:8080/upload', formData, {
            headers: {
            'Content-Type': 'multipart/form-data',
            "Authorization":token,
            },
        });
        
        setUploadFileData(res)
        setModalContent({
            header: "Upload Success",
            message: "Success upload file",
        });
        if (modalRef.current) {
            modalRef.current.openModal();
        }

        console.log(res)

        getReport()
        setHasUploafFile(true)
        localStorage.setItem("id_report",res.data.id_report);
        localStorage.setItem("date",res.data.date);
        setLoading(false)
        } catch (error) {
            console.error('Error uploading file:', error);
            setModalContent({
                header:"Upload Failed",
                message:"Error uploading file",
            })
            if (modalRef.current) {
                modalRef.current.openModal();
            }
            setLoading(false)
        }
    };

    const Logout = ()=>{
       
        localStorage.clear()
        navigate("/")
    }

    const getReport = async ()=>{
        
        try{
            const response = await axios.get("http://localhost:8080/get/report/"+id_user,
                {headers:{
                    'Content-Type':'application/json',
                    "Authorization":token,
                }}
            )
            
            const data = response.data
            setReport(data)


        }catch(error){
            console.log(error)
        }
    }

    const deleteReport = async (id)=>{
        try{
            const response = await axios.delete("http://localhost:8080/delete/report/"+id,
                {headers:{
                    'Content-Type':'application/json',
                    "Authorization":token,
                }}
            )
            
            getReport()
        }catch(error){
            console.log(error)
        }
    }

    const aiChat = async ()=>{
        try{
            
            setIsLoadingChat(true)
          
            const formData = new FormData();
            formData.append("id_user", localStorage.getItem("id_user"));
            formData.append("id_report",localStorage.getItem("id_report"));
            formData.append("date",localStorage.getItem("date"));
            formData.append("query",question);
            formData.append("document", uploadFileData.data.stringText);

            setQuestion("")

            const response = await axios.post("http://localhost:8080/chat",formData,
                {
                    headers:{
                        'Content-Type': 'multipart/form-data',
                        "Authorization":token,
                    }
                }
            )
            console.log(response)

            

            setAIResponse((prevState) => ({
                ...prevState, 
                question: [...prevState.question, response.data.Question], 
                answer: [...prevState.answer, response.data.responseAI],  
            }));

            setIsLoadingChat(false)
            
        }
        catch(error){
            console.log(error)
            setIsLoadingChat(false)
        }
    }
    
    
    return(
        <>
            <Modal ref={modalRef} header={modalContent.header} message={modalContent.message} />
            <div  className='absolute right-0 top-0 m-4 mb-8  '>
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
                            <h1 className='text-xl text-center bg-slate-600 text-white rounded-lg mb-5' onClick={reset} >New Chat</h1>
                            {/* Sidebar content here */}                          
                            
                            {
                                
                            reports && reports.map(item=>(
                                <div className='flex w-full justify-between items-center ' key={item.id} >
                                    
                                    <Link to={"/Chat/"+item.id}
                                        className='w-3/4 mb-2  text-black rounded-xl p-2'
                                        
                                    >
                                        {item.date}
                                    </Link>
                                    <DeleteIcon onClick={()=>deleteReport(item.id)} style={{ fill: 'red' }} />
                                </div>
                            ))}
                            
                            
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
                        onClick={handleUpload}
                    >
                        {loading ? <Loader/>:"Upload"}
                        
                    </button>

                </div>
                <div
                    ref={chatContainerRef} // Pasang referensi ke elemen ini
                    className="chat-container max-h-screen overflow-y-auto"
                    style={{ maxHeight: "80vh" }} // Atur tinggi maksimal sesuai kebutuhan
                >
                    {aiResponse.question.map((q, index) => (
                        <React.Fragment key={index}>
                            <ChatEnd message={q} />
                            <ChatStart message={aiResponse.answer[index]} />
                        </React.Fragment>
                    ))}
                </div>
            </div>
            

            <div className="fixed bottom-0 left-0 right-0 w-full p-4 bg-white shadow-lg  flex items-center space-x-2 h-10">
                <input 
                    type="text" 
                    value={question} 
                    onChange={(e)=>{setQuestion(e.target.value)}}
                    className="input input-bordered  w-full mb-14 p-4" 
                />
                <button >
                {hasUploadFile === true && (
                    isLoadingChat? (
                        <Loader /> 
                    ) : (
                        <Send 
                            className='Send mb-10' 
                            onClick={aiChat} 
                            width={40} 
                            height={40} 
                        />
                    )
                )}
                    
                </button>
            </div>
        </>
    )
}

export default ChatPage;