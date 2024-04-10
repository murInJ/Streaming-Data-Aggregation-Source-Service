package encoder

import (
	"errors"
	"fmt"
	"image"
	"unsafe"
)

// #cgo pkg-config: libavcodec libavutil libswscale
// #include <libavcodec/avcodec.h>
// #include <libavutil/imgutils.h>
// #include <libswscale/swscale.h>
import "C"

// H264Encoder is a wrapper around FFmpeg's H264 encoder.
type H264Encoder struct {
	codecCtx *C.AVCodecContext
	srcFrame *C.AVFrame
	dstFrame *C.AVFrame
	swsCtx   *C.struct_SwsContext
}

// Initialize initializes a H264Encoder.
func (e *H264Encoder) Initialize(width, height int) error {
	codec := C.avcodec_find_encoder(C.AV_CODEC_ID_H264)
	if codec == nil {
		return fmt.Errorf("avcodec_find_encoder() failed")
	}

	e.codecCtx = C.avcodec_alloc_context3(codec)
	if e.codecCtx == nil {
		return fmt.Errorf("avcodec_alloc_context3() failed")
	}

	e.codecCtx.bit_rate = 400000 // 设置比特率
	e.codecCtx.width = C.int(width)
	e.codecCtx.height = C.int(height)
	e.codecCtx.time_base.num = 1 // 设置时间基准
	e.codecCtx.time_base.den = 25
	e.codecCtx.gop_size = 10 // 设置关键帧间隔
	e.codecCtx.max_b_frames = 1
	e.codecCtx.pix_fmt = C.AV_PIX_FMT_YUV420P // 使用 FFmpeg 提供的常量

	res := C.avcodec_open2(e.codecCtx, codec, nil)
	if res < 0 {
		C.avcodec_close(e.codecCtx)
		return fmt.Errorf("avcodec_open2() failed")
	}

	e.srcFrame = C.av_frame_alloc()
	if e.srcFrame == nil {
		C.avcodec_close(e.codecCtx)
		return fmt.Errorf("av_frame_alloc() failed")
	}

	e.srcFrame.format = C.int(int(e.codecCtx.pix_fmt))
	e.srcFrame.width = e.codecCtx.width
	e.srcFrame.height = e.codecCtx.height
	res = C.av_frame_get_buffer(e.srcFrame, 0)
	if res < 0 {
		return fmt.Errorf("av_frame_get_buffer() failed")
	}

	e.dstFrame = C.av_frame_alloc()
	if e.dstFrame == nil {
		C.avcodec_close(e.codecCtx)
		return fmt.Errorf("av_frame_alloc() failed")
	}

	e.dstFrame.format = C.AV_PIX_FMT_YUV420P
	e.dstFrame.width = e.codecCtx.width
	e.dstFrame.height = e.codecCtx.height
	res = C.av_frame_get_buffer(e.dstFrame, 0)
	if res < 0 {
		return fmt.Errorf("av_frame_get_buffer() failed")
	}

	e.swsCtx = C.sws_getContext(e.codecCtx.width, e.codecCtx.height, C.AV_PIX_FMT_RGBA,
		e.codecCtx.width, e.codecCtx.height, e.codecCtx.pix_fmt, C.SWS_BICUBIC, nil, nil, nil)
	if e.swsCtx == nil {
		return fmt.Errorf("sws_getContext() failed")
	}

	return nil
}

// Close closes the encoder.
func (e *H264Encoder) Close() {
	C.sws_freeContext(e.swsCtx)
	C.av_frame_free(&e.srcFrame)
	C.av_frame_free(&e.dstFrame)
	C.avcodec_close(e.codecCtx)
}

// Encode encodes an image to H.264.
func (e *H264Encoder) Encode(img image.Image) ([]byte, error) {
	// Convert image to RGBA if not already RGBA
	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		return nil, errors.New("image must be RGBA")
	}

	// Convert RGBA image to YUV420P
	inImg := (*C.uint8_t)(unsafe.Pointer(&rgbaImg.Pix[0]))
	inLinesize := (C.int)(rgbaImg.Stride)
	outImg := (**C.uint8_t)(unsafe.Pointer(&e.srcFrame.data[0]))
	outLinesize := (*C.int)(unsafe.Pointer(&e.srcFrame.linesize[0]))
	C.sws_scale(e.swsCtx, &inImg, &inLinesize, 0, (C.int)(rgbaImg.Rect.Dy()), outImg, outLinesize)

	// Encode frame
	res := C.avcodec_send_frame(e.codecCtx, e.srcFrame)
	if res < 0 {
		return nil, fmt.Errorf("avcodec_send_frame() failed")
	}

	var pkt C.AVPacket
	C.av_init_packet(&pkt)
	defer C.av_packet_unref(&pkt)

	res = C.avcodec_receive_packet(e.codecCtx, &pkt)
	if res < 0 {
		return nil, fmt.Errorf("avcodec_receive_packet() failed")
	}

	return C.GoBytes(unsafe.Pointer(pkt.data), C.int(pkt.size)), nil
}
