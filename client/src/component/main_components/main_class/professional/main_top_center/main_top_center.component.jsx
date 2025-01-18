import prof1 from '../../../../../assets/prp1.jpeg'
import prof2 from '../../../../../assets/prp3.jpeg'
import prof3 from '../../../../../assets/prp4.jpeg'
import prof4 from '../../../../../assets/prp5.jpeg'
import './main_top_center.component.css'

const Professional_top_center = () => {
    return (
        <div className='main_prof'>
            <div className="main_top_center_prof main_top_center">

                <div className="main_top_right_prof">
                    <div className="main_top_right_prof_img">
                        <img src={prof1} alt="" />
                    </div>
                    <div className="main_top_right_prof_text">
                        <p>
                            <span>В</span> техникуме на отделении заочной формы обучения реализуются программы
                            дополнительного профессионального образования по специальностям:
                        </p>
                    </div>

                </div>
                <div className="main_top_right_prof">
                    <div className="main_top_right_prof_img">
                        <img src={prof2} alt="" />
                    </div>
                    <div className="main_top_right_prof_text" id="bg">
                        <p>
                            <span>Программы</span> профессиональной переподготовки направлены на получение
                            компетенций,
                            необходимых для выполнения нового вида профессиональной деятельности.
                        </p>
                    </div>

                </div>

            </div>
            <div className="main_top_center_prof main_top_center">
                <div className="main_top_right_prof" >
                    <div className="main_top_right_prof_img">
                        <img src={prof3} alt="" />
                    </div>
                    <div className="main_top_right_prof_text" id="bg">
                        <p>
                            <span>C</span> 2008 по 2020 г. получили дипломы профессиональной переподготовки 862
                            специалиста, занятых
                            в нефтяной, строительной отрасли и энергетическом комплексе республики.
                        </p>
                    </div>

                </div>
                <div className="main_top_right_prof" >
                    <div className="main_top_right_prof_img">
                        <img src={prof4} alt="" />
                    </div>
                    <div className="main_top_right_prof_text">
                        <p>
                            <span>К</span> освоению дополнительных профессиональных образовательных программ
                            допускаются:
                        </p>
                        <ul>
                            <li>лица, имеющие среднее профессиональное и (или) высшее образование;</li>
                            <li>лица, получающие среднее профессиональное и (или) высшее образование.</li>
                        </ul>
                    </div>

                </div>
            </div>
        </div>
    )
}

export default Professional_top_center;