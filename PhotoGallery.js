import React from 'react';
import Gallery from 'react-grid-gallery';
import PropTypes from 'prop-types';
import { useState, useEffect } from "react";
import{DeleteOutlined} from "@ant-design/icons";
import{BASE_URL, TOKEN_KEY} from "../constants";
import axios from "axios";
import {message, Button} from "antd";


const captionStyle = {
    backgroundColor: "rgba(0, 0, 0, 0.6)",
    maxHeight: "240px",
    overflow: "hidden",
    position: "absolute",
    bottom: "0",
    width: "100%",
    color: "white",
    padding: "2px",
    fontSize: "90%"
};

const wrapperStyle = {
    display: "block",
    minHeight: "1px",
    width: "100%",
    border: "1px solid #ddd",
    overflow: "auto"
};

//overlay: style + 属性
function PhotoGallery(props) {
    const [images, setImages] = useState(props.images);
    const [curImgIdx, setCurImgIdx] = useState(0);

    const imageArr = images.map( image => {
        return {
            ...image,
            customOverlay: (
                <div style={captionStyle}>
                    <div>{`${image.user}: ${image.caption}`}</div>
                </div>
            )
        }
    });
    const onDeleteImage = () => {
        if (window.confirm(`Are you sure you want to delete this image?`)){
            //find cur image, delete from image array
            //newimage: store all images without index = curimgindex
            const curImg = images[curImgIdx];
            const newImageArr = images.filter((img, index) => index !== curImgIdx);
            console.log('delete image ', newImageArr);
            //inform server to delete the image
            const opt = {
                method: 'DELETE',
                url: `${BASE_URL}/post/${curImg.postId}`,
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem(TOKEN_KEY)}`
                }
            };

            //if deleted, set image states: setImages
            axios(opt)
                .then( res => {
                    console.log('delete result -> ', res);
                    // case1: success
                    if(res.status === 200) {
                        // step1: set state
                        setImages(newImageArr);
                    }
                })
                .catch( err => {
                    // case2: fail
                    message.error('Fetch posts failed!');
                    console.log('fetch posts failed: ', err.message);
                })
            //setImages(newImageArr);
        }
    }

    //find cur image index
    //index is from Gallery
    const onCurrentImageChange = index => {
        console.log('curIdx ', index);
        setCurImgIdx(index) //record index to states
    }

    //when props changes, update states.
    useEffect(() => {
        setImages(props.images);
    }, [props.images])

    return (
        <div style={wrapperStyle}>
            <Gallery
                images={imageArr}
                enableImageSelection={false}
                backdropClosesModal={true}
                currentImageWillChange={onCurrentImageChange} //current image
                customControls={[
                    <Button style={{marginTop: "10px", marginLeft: "5px"}}
                            key="deleteImage"
                            type="primary"
                            icon={<DeleteOutlined />}
                            size="small"
                            onClick={onDeleteImage}
                    >Delete Image</Button>
                ]}
            />
        </div>
    );
}

//类型校验
PhotoGallery.propTypes = {
    images: PropTypes.arrayOf(
        PropTypes.shape({
            user: PropTypes.string.isRequired,
            caption: PropTypes.string.isRequired,
            src: PropTypes.string.isRequired,
            thumbnail: PropTypes.string.isRequired,
            thumbnailWidth: PropTypes.number.isRequired,
            thumbnailHeight: PropTypes.number.isRequired,
        })
    ).isRequired
};

export default PhotoGallery;
