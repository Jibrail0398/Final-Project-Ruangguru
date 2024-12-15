import { Routes,Route } from "react-router-dom";
import RegisterPage from "./RegisterPage";
import LoginPage from "./LoginPage";
import ChatPage from "./ChatPage";
import ProtectedRoute from "./ProtectedRoutes";
import ChatHistory from "./ChatHistory";

function MainRoutes (){

    return(
        <>
            <Routes>
                <Route path="/" element={<LoginPage/>} ></Route>
                <Route path="/Register" element={<RegisterPage/>} ></Route>
                <Route 
                    path="/Chat"
                
                >
                    <Route index element={
                        <ProtectedRoute>
                            <ChatPage/>
                        </ProtectedRoute>
                        
                        } ></Route>
                    <Route path=":id" element={
                        <ProtectedRoute>
                            <ChatHistory/>
                        </ProtectedRoute>
                        
                        } ></Route>
                </Route>
            </Routes>
        
        </>
    )
}

export default MainRoutes;