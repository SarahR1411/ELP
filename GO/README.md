### **README: Photo Restoration Project**

---

#### **Overview**
This project restores old or damaged photographs using a combination of noise reduction, scratch/damage repair, and color correction techniques. Implemented in Go, the program processes an input image to reduce noise, repair scratches or stains, and enhance its overall visual quality.

---

#### **How It Works**
1. **Load and Save Images:**
   - The program reads images from the specified file path and saves processed images in JPEG format.

2. **Noise Reduction:**
   - Uses Gaussian-like filtering to reduce noise by averaging neighboring pixel colors. This smooths the image while retaining key features.

3. **Scratch and Damage Repair:**
   - Detects high-contrast, bright pixels (potential scratches or stains) by thresholding.
   - Repairs the damaged areas by replacing them with the median color of surrounding pixels, effectively blending them into the background.

4. **Color Correction:**
   - Implements **Histogram Equalization** to enhance image contrast. It redistributes pixel intensities evenly across the image, improving brightness and visual quality.

---

#### **Code Details**
1. **Key Functions:**
   - `loadImage(imagePath string)`: Reads an image from a file.
   - `saveImage(img image.Image, outputPath string)`: Writes an image to a file in JPEG format.
   - `detectWhiteLinesAndStains(img image.Image)`: Identifies potential scratches or stains using brightness thresholds.
   - `applySmoothing(img image.Image)`: Reduces noise and smooths the image using median filtering.
   - `HistEqual(img image.Image)`: Performs histogram equalization to adjust image contrast.
   - `getMedianColor(img image.Image, x, y int)`: Calculates the median color of surrounding pixels for restoration.

2. **Workflow:**
   - **Step 1:** The image is loaded.
   - **Step 2:** Scratches and stains are detected using a brightness threshold.
   - **Step 3:** The image undergoes histogram equalization for contrast improvement.
   - **Step 4:** Noise reduction and smoothing are applied to refine the final image.
   - **Step 5:** The restored image is saved.

---

#### **Filters Explanation**
1. **Noise Reduction (Gaussian-like Smoothing):**
   - Replaces each pixel with the average color of its neighbors. This reduces noise by blending pixel values, resulting in a smoother appearance.

2. **Scratch Repair (Median Filtering):**
   - Bright pixels, indicative of damage, are replaced with the median color of surrounding pixels. This method preserves edges and avoids excessive blurring.

3. **Color Correction (Histogram Equalization):**
   - Balances the intensity distribution of pixels across the image. Dark areas become brighter, and contrast improves, making the image more visually appealing.

---

#### **Usage**
1. Set the image path:  
   Update `imagePath` with the file path to your input image.
   
2. Run the program:  
   Execute the `main` function to process the image.

3. Output:  
   The restored image will be saved as `restored_photo.jpg`.

---

#### **Future Improvements**
- Add advanced noise reduction techniques (e.g., bilateral filtering).
- Incorporate machine learning for more robust scratch detection and restoration.
- Improve resolution enhancement using modern upscaling methods. 

Enjoy restoring your photos! 📸
