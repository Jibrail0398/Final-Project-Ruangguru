import { ReactComponent as MyIcon } from '../assets/visionary.svg';
import { Link } from 'react-router-dom';

function RegisterPage(){
    return(
        <>
            <div className="flex justify-center items-center h-screen " >

                <div className="card bg-base-100 w-full h-full md:w-3/4 md:h-1/2 lg:h-full shadow-xl">
                    <div className="card-body  lg:flex-row ">

                        <div className=' lg:w-1/2' >

                            <h2 className="card-title text-3xl ">Register!</h2>
                            <label className="form-control w-full max-w-xs">
                                <div className="label">
                                    <span className="label-text text-lg">Username</span>
                                </div>
                                <input type="text" placeholder="Type here" className="input input-bordered w-full max-w-xs" />
                            </label>
                            <label className="form-control w-full max-w-xs">
                                <div className="label">
                                    <span className="label-text text-lg">Email</span>
                                </div>
                                <input type="email" placeholder="Type here" className="input input-bordered w-full max-w-xs" />
                            </label>
                            <label className="form-control w-full max-w-xs">
                                <div className="label">
                                    <span className="label-text text-lg">Password</span>
                                </div>
                                <input type="password" placeholder="Type here" className="input input-bordered w-full max-w-xs" />
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

                        <button className="btn btn-success w-3/4 mx-auto m-2 text-white ">Register</button>
                    </div>
                </div>
            </div>
        </>
    )
}

export default RegisterPage;