import { Routes,Route } from "react-router-dom";
import RegisterPage from "./RegisterPage";
import LoginPage from "./LoginPage";
import ChatPage from "./ChatPage";


function MainRoutes (){

    return(
        <>
            <Routes>
                <Route path="/" element={<LoginPage/>} ></Route>
                <Route path="/Register" element={<RegisterPage/>} ></Route>
                <Route path="/Chat">
                    <Route index element={<ChatPage/>} ></Route>
                    <Route></Route>
                </Route>
            </Routes>
        
        </>
    )
}

export default MainRoutes;