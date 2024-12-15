function ChatStart(props){
    return(
        <div className="chat chat-start">
            <div className="chat-bubble chat-bubble-success text-white ml-4 mt-4">
                {props.message}
            </div>
        </div>

    )
}

export default ChatStart;