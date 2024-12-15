import React, { forwardRef, useImperativeHandle, useRef } from "react";

const Modal = forwardRef((props, ref) => {
    const modalRef = useRef();
  
    function openModal() {
      const myModal = modalRef.current;
      if (myModal) {
        myModal.showModal();
      }
    }
  
    // Expose the openModal method to the parent
    useImperativeHandle(ref, () => ({
      openModal,
    }));
  
    return (
      <>
        <dialog id="my_modal_1" ref={modalRef} className="modal">
          <div className="modal-box">
            <h3 className="font-bold text-lg">{props.header}</h3>
            <p className="py-4">{props.message}</p>
            <div className="modal-action">
              <form method="dialog">
                <button className="btn">Ok</button>
              </form>
            </div>
          </div>
        </dialog>
      </>
    );
  });

export default Modal;