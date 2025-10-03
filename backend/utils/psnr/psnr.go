package psnr

import (
    "fmt"
    "math"
)





func CalculatePSNR(original, stego []byte) (float64, error) {
    if len(original) != len(stego) {
        return 0, fmt.Errorf("audio lengths don't match: original=%d, stego=%d", len(original), len(stego))
    }
    
    if len(original) == 0 {
        return 0, fmt.Errorf("empty audio data")
    }
    
    
    if len(original)%2 != 0 {
        return 0, fmt.Errorf("audio data length must be even for 16-bit PCM")
    }
    
    numSamples := len(original) / 2
    
    
    var mse float64
    for i := 0; i < numSamples; i++ {
        
        originalSample := int16(original[i*2]) | int16(original[i*2+1])<<8
        stegoSample := int16(stego[i*2]) | int16(stego[i*2+1])<<8
        
        
        diff := float64(originalSample) - float64(stegoSample)
        mse += diff * diff
    }
    mse /= float64(numSamples)
    
    
    if mse == 0 {
        return math.Inf(1), nil
    }
    
    
    maxValue := 32767.0
    psnr := 10 * math.Log10((maxValue*maxValue)/mse)
    
    return psnr, nil
}


func CalculatePSNRFloat(original, stego []float64) (float64, error) {
    if len(original) != len(stego) {
        return 0, fmt.Errorf("audio lengths don't match: original=%d, stego=%d", len(original), len(stego))
    }
    
    if len(original) == 0 {
        return 0, fmt.Errorf("empty audio data")
    }
    
    
    var mse float64
    for i := 0; i < len(original); i++ {
        diff := original[i] - stego[i]
        mse += diff * diff
    }
    mse /= float64(len(original))
    
    
    if mse == 0 {
        return math.Inf(1), nil
    }
    
    
    maxValue := 1.0
    psnr := 10 * math.Log10((maxValue*maxValue)/mse)
    
    return psnr, nil
}


func CalculatePSNRByte(original, stego []byte) (float64, error) {
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


func DetectAudioFormat(original, stego []byte) (float64, string, error) {
    
    if len(original)%2 == 0 && len(original) > 0 {
        psnr16, err := CalculatePSNR(original, stego)
        if err == nil {
            return psnr16, "16-bit PCM", nil
        }
    }
    
    
    psnr8, err := CalculatePSNRByte(original, stego)
    if err != nil {
        return 0, "", err
    }
    
    return psnr8, "8-bit PCM", nil
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
        return "Good Quality (Acceptable)"
    case psnr >= 20:
        return "Fair Quality"
    default:
        return "Poor Quality (Damaged)"
    }
}
func FormatPSNR(psnr float64) string {
    if math.IsInf(psnr, 1) {
        return "∞ dB"
    }
    return fmt.Sprintf("%.2f dB", psnr)
}

func CalculateMSE(original, stego []byte) (float64, error) {
    if len(original) != len(stego) {
        return 0, fmt.Errorf("audio lengths don't match")
    }
    
    if len(original) == 0 {
        return 0, fmt.Errorf("empty audio data")
    }
    
    var mse float64
    numSamples := len(original) / 2 
    
    for i := 0; i < numSamples; i++ {
        originalSample := int16(original[i*2]) | int16(original[i*2+1])<<8
        stegoSample := int16(stego[i*2]) | int16(stego[i*2+1])<<8
        
        diff := float64(originalSample) - float64(stegoSample)
        mse += diff * diff
    }
    
    return mse / float64(numSamples), nil
}