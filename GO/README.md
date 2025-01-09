### Project README: **Image Processing Toolkit**

---

#### **Overview**
This project focuses on building a toolkit for basic image processing tasks, such as noise reduction, color correction, resolution enhancement, and scratch/damage repair. These tools aim to improve image quality and fix common visual defects with simple and practical approaches.

---

#### **Features to Implement**
1. **Noise Reduction:**
   - Implement a box filter (blur) by averaging RGB values of neighboring pixels.
   - Alternatively, use Gaussian filtering for better smoothing and noise reduction.

2. **Color Correction:**
   - Methods to include:
     - Histogram equalization.
     - White balance adjustment.
     - Channel-wise intensity scaling.

3. **Resolution Enhancement:**
   - Basic upscaling methods:
     - Nearest-neighbor interpolation.
     - Bilinear interpolation.
   - Note: These methods are simple but produce subpar results compared to advanced machine learning-based techniques.

4. **Scratch/Damage Repair:**
   - Detect bright pixels with high contrast to their surroundings using edge detection algorithms.
   - Replace damaged pixels with the average values of neighboring pixels for seamless repair.

---

#### **Development Plan**
1. **Priority Tasks:**
   - **Noise Reduction:** Start with implementing Gaussian filtering for smoother and cleaner images.
   - **Scratch/Damage Repair:** Develop algorithms to detect and repair scratches or bright lines.

2. **Secondary Tasks:**
   - **Color Correction:** Focus on enhancing visual realism and naturalness by correcting color balance.

3. **Final Task:**
   - **Resolution Enhancement:** Add basic upscaling methods for increasing image size with minimal complexity.

---

#### **Progress So Far**
- **Core Functions:**
  - `LoadImage` and `SaveImage` functions for reading and saving images.
  - `ScratchDetection` function for identifying scratches and bright lines.

- **Noise Reduction:**
  - `ApplySmoothing`: Implements Gaussian filtering for reducing noise.
  - `GetMedianColor`: Computes the average color of surrounding pixels, used for smoothing.

---

#### **Next Steps**
1. Fine-tune the scratch detection and repair process.
2. Complete the noise reduction module by adding box filtering.
3. Implement the color correction methods (histogram equalization, white balance).
4. Develop basic resolution enhancement methods.

---

#### **Notes**
- This project avoids machine learning-based methods to keep implementation simple and focus on classic image processing techniques.
- The workflow prioritizes repair and enhancement for practical and visually noticeable improvements.
