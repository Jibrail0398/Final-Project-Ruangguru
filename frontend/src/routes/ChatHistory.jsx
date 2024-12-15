import {ReactComponent as ListIcon } from '../assets/list.svg';
import {ReactComponent as DeleteIcon} from '../assets/trash-fill.svg'
import { useLocation, Link,useNavigate,useParams } from 'react-router-dom';
import React, { useState,useEffect,useRef } from "react";
import ChatStart from "../component/ChatStart";
import ChatEnd from "../component/ChatEnd";
import axios from "axios";
import  {ReactComponent as Send} from "../assets/send.svg";
import Modal from "../component/Modal"
import Loader from "../component/Loader"

function ChatHistory(){
    //Menangkap id dari url parameter
    let { id } = useParams();
    //untuk navigate
    let navigate = useNavigate();
    //referensi modal untuk keperluan membuka modal dari elemen parent
    const modalRef = useRef();
    
    //untuk menangkap data report untuk sidebar
    const [reports,setReport] = useState(null);

    const [isNewChat,setNewChat] = useState(false);

    //pertanyaan di kolom chat AI yang dibawah
    const [question,setQuestion] = useState("");
    // Referensi untuk container chat
    const chatContainerRef = useRef(null);
    //untuk menangkap response dari AI
    const [aiResponse,setAIResponse] = useState({
        question:[],
        answer:[],
    });
    //loading saat chat
    const [isLoadingChat,setIsLoadingChat] = useState(false)
    
    const location = useLocation();
    

    const isActive = (path) => location.pathname === path;



    //ini untuk modal
    const [modalContent, setModalContent] = useState({
        header: "",
        message: "",
    });


    //Ambil data id_user dan token dari localStorage
    const id_user = localStorage.getItem("id_user")
    const token = localStorage.getItem("token")
    
    //fungsi untuk melakukan fitur new chat
    const reset = ()=>{
        setQuestion("")
        
    }


    //react lifecycle untuk menangkap data report yang ada di sidebar
    useEffect(()=>{
        getReport() 
        getHistoryChat(id)

        return()=>{
            localStorage.removeItem("date")
            localStorage.removeItem("stringText")
            localStorage.removeItem("id_report")
        }
    },[id])

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
    
    
    //fungsi ambil history chat

    const getHistoryChat = async (id) =>{
        try{
            const response = await axios.get('http://localhost:8080/chat/get/'+id,{
                headers:{
                    "Content-Type":"application/json",
                    "Authorization":token,
                }
            });
            
             // Ekstrak data dari respons API
            const chatData = response.data; // Mengakses array "data"
            
            // Memisahkan pertanyaan (question) dan jawaban (response)
            const questions = chatData.map(item => item.question);
            const answers = chatData.map(item => item.response);

            // Update state dengan data yang didapat
            setAIResponse({
                question: questions,
                answer: answers,
            });
            
        }catch(error){
            console.log(error)

        }
    }


    //untuk melakukan logout
    const Logout = ()=>{
       
        localStorage.clear()
        navigate("/")
    }

    //untuk mendapatkan data report yang ada di sidebar
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
            data.forEach(element => {
                if (element.id == id){
                    
                    localStorage.setItem("date",element.date)
                    localStorage.setItem("stringText",element.stringText)
                    localStorage.setItem("id_report",element.id)

                }
            });


        }catch(error){
            console.log(error)
        }
    }

    //Untuk melakukan deleteReport yang ada di sidebar
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

    //Melakukan komunikasi dengan AI
    const aiChat = async ()=>{
        try{
            
            setIsLoadingChat(true)
          
            const formData = new FormData();
            formData.append("id_user", localStorage.getItem("id_user"));
            formData.append("id_report",localStorage.getItem("id_report"));
            formData.append("date",localStorage.getItem("date"));
            formData.append("query",question);
            formData.append("document", localStorage.getItem("stringText"));

            setQuestion("")

            const response = await axios.post("http://localhost:8080/chat",formData,
                {
                    headers:{
                        'Content-Type': 'multipart/form-data',
                        "Authorization":token,
                    }
                }
            )
            
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
                            <h1 className='text-xl text-center cursor-pointer  bg-slate-600 text-white rounded-lg mb-5' onClick={()=>navigate("/chat")} >New Chat</h1>
                            {/* Sidebar content here */}                          
                            
                            {
                                
                            reports && reports.map(item=>(
                                <div className='flex w-full justify-between items-center ' key={item.id} >
                                    
                                    <Link to={"/Chat/"+item.id}
                                    
                                        className={` w-3/4 text-lg mb-2 ${isActive('/Chat/'+item.id) ? 'mb-2 bg-slate-400 text-white rounded-xl p-2' : 'w-3/4 mb-2  text-black rounded-xl p-2'}`}
                                        
                                    >
                                        {item.date}
                                    </Link>
                                    <DeleteIcon onClick={()=>deleteReport(item.id)} style={{ fill: 'red' }} />
                                </div>
                            ))}
                            
                            
                        </ul>
                    </div>
                    
                </div>
            
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

            <div className="fixed bottom-0 left-0 right-0 w-full p-4 bg-white shadow-lg  flex items-center space-x-2 h-10">
                <input 
                    type="text" 
                    value={question} 
                    onChange={(e)=>{setQuestion(e.target.value)}}
                    className="input input-bordered  w-full mb-14 p-4" 
                />
                <button >
                {isLoadingChat? (
                        <Loader /> 
                    ) : (
                        <Send 
                            className='Send mb-10' 
                            onClick={aiChat} 
                            width={40} 
                            height={40} 
                        />
                )}
                    
                </button>
            </div>
        </>
    )
}

export default ChatHistory;