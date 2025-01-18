import teacher1 from "../../../../../assets/teacher_1.jpeg";
import teacher2 from "../../../../../assets/teacher_2.jpeg";
import teacher3 from "../../../../../assets/teacher_3.jpeg";
import teacher4 from "../../../../../assets/teacher_4.jpeg";

const Teacher_top_center = () => {
    return (
        <>
            <div className="main_top_center_prof main_top_center">

                <div className="main_top_right_prof">
                    <div className="main_top_right_prof_img">
                        <img src={teacher1} alt=""/>
                    </div>
                    <div className="main_top_right_prof_text main_top_right_teacher_text">
                        <p>
                            Валова Юлия Александровна
                        </p>
                        <p>
                            Профессиональная переподготовка  <br />
                            Транспортное направление  <br />
                            Курсы подготовки и переподготовки водителей и специалистов занятых перевозкой опасных
                            грузов
                        </p>
                        <p>
                            Кабинет записи:  <br />
                            2 этаж,  <br />
                            221 кабинет  <br />
                            т. 8(8553) 33-06-37
                        </p>
                    </div>

                </div>
                <div className="main_top_right_prof">
                    <div className="main_top_right_prof_img">
                        <img src={teacher2} alt=""/>
                    </div>
                    <div className="main_top_right_prof_text main_top_right_teacher_text" id="bg">
                        <p>
                            Терентьева Елена Вячеславовна
                        </p>
                        <p>Профессиональная подготовка по рабочим профессиям</p>
                        <p>Кабинет записи:  <br />
                            2 этаж,  <br />
                            208 кабинет  <br />
                            т. 8(8553) 33-48-60 (экономический отдел)</p>
                    </div>

                </div>

            </div>
            <div className="main_top_center_prof main_top_center">
                <div className="main_top_right_prof">
                    <div className="main_top_right_prof_img">
                        <img src={teacher3} alt=""/>
                    </div>
                    <div className="main_top_right_prof_text main_top_right_teacher_text" id="bg">
                        <p>
                            Алаев Владимир Петрович
                        </p>
                        <p>Автошкола</p>
                        <p>Кабинет записи:  <br />
                            2 этаж,  <br />
                            211 кабинет  <br />
                            т. 8(8553) 33-48-58</p>
                    </div>

                </div>
                <div className="main_top_right_prof">
                    <div className="main_top_right_prof_img">
                        <img src={teacher4} alt=""/>
                    </div>
                    <div className="main_top_right_prof_text main_top_right_teacher_text">
                        <p>
                            Осянина Лилия Раисовна
                        </p>
                        <p>Курсы повышения квалификации по линии СРО</p>
                        <p>Кабинет записи:  <br />
                            2 этаж,  <br />
                            208 кабинет  <br />
                            т. 8(8553) 33-48-60</p>
                    </div>

                </div>
            </div>


        </>
    )
};

export default Teacher_top_center;