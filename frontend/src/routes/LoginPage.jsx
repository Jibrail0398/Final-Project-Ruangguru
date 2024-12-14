import { ReactComponent as MyIcon } from '../assets/visionary.svg';
import { Link } from 'react-router-dom';

function LoginPage(){
    return(
        <>
            <div className="flex justify-center items-center h-screen " >

                <div className="card bg-base-100 w-full h-full md:w-3/4 md:h-1/2 lg:h-full shadow-xl">
                    <div className="card-body  lg:flex-row ">

                        <div className=' lg:w-1/2' >

                            <h2 className="card-title text-3xl ">Login!</h2>
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
                    <div className="card-actions justify-end">
                        <button className="btn btn-primary w-3/4 mx-auto text-white ">Login</button>
                        <Link to="/Register" className="btn btn-success w-3/4 mx-auto mb-4 text-white " >Register</Link>
                        
                    </div>
                </div>
            </div>
        </>
    )
}

export default LoginPage;