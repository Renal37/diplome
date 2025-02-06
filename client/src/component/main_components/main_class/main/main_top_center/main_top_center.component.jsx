import React from "react";
import Slider from "react-slick";
import "slick-carousel/slick/slick.css";    
import "slick-carousel/slick/slick-theme.css";
import "./main_top_center.component.css";
import apt from "../../../../../assets/apt.png";

const Main_top_center = () => {
    const settings = {
        dots: true,
        infinite: true,
        speed: 500,
        slidesToShow: 1,
        slidesToScroll: 1,
        autoplay: true,
        autoplaySpeed: 3000,
    };

    return (
        <div className="main_top_center_slider">
            {/* Статичное изображение */}
            <img src={apt} alt="Наш техникум" className="static-image" />

            {/* Слайдер с текстовым содержимым */}
            <Slider {...settings} className="main_top_slider">
                {/* Слайд 1 */}
                <div className="slider-item">
                    <div className="slider-content">
                        <h3>Почему наш сайт?</h3>
                        <p>
                            Мы предлагаем качественные образовательные программы, разработанные опытными преподавателями.
                        </p>
                    </div>
                </div>

                {/* Слайд 2 */}
                <div className="slider-item">
                    <div className="slider-content">
                        <h3>Почему именно у нас?</h3>
                        <p>
                            Наши курсы помогут вам получить актуальные знания и навыки, необходимые для успешной карьеры.
                        </p>
                    </div>
                </div>

                {/* Слайд 3 */}
                <div className="slider-item">
                    <div className="slider-content">
                        <h3>Преимущества наших курсов</h3>
                        <p>
                            Гибкий график обучения, доступные цены и поддержка на всех этапах обучения.
                        </p>
                    </div>
                </div>
            </Slider>
        </div>
    );
};

export default Main_top_center;