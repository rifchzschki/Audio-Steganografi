package psnr

import(
	"fmt"
	"math"
)

func CalculatePSNR(original, stego []byte) (float64, error) {
    if len(original) != len(stego) {
        return 0, fmt.Errorf("audio lengths don't match")
    }
    
    if len(original) == 0 {
        return 0, fmt.Errorf("empty audio data")
    }
    
    var mse float64
    for i := 0; i < len(original); i++ {
        diff := float64(original[i]) - float64(stego[i])
        mse += diff * diff
    }
    mse /= float64(len(original))
    
    if mse == 0 {
        return math.Inf(1), nil
    }
    maxValue := 255.0
    psnr := 10 * math.Log10((maxValue*maxValue)/mse)
    
    return psnr, nil
}

func IsQualityAcceptable(psnr float64) bool {
    if math.IsInf(psnr, 1) {
        return true
    }
    return psnr >= 30.0
}

func GetQualityStatus(psnr float64) string {
    if math.IsInf(psnr, 1) {
        return "Perfect Quality (Identical)"
    }
    
    switch {
    case psnr >= 50:
        return "Excellent Quality"
    case psnr >= 40:
        return "Very Good Quality"
    case psnr >= 30:
        return "Good Quality"
    case psnr >= 20:
        return "Fair Quality"
    default:
        return "Poor Quality"
    }
}

func FormatPSNR(psnr float64) string {
    if math.IsInf(psnr, 1) {
        return "∞ dB"
    }
    return fmt.Sprintf("%.2f dB", psnr)
}



