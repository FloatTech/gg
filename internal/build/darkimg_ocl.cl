kernel void isdark(
    read_only image2d_t inputImg,
    sampler_t smp,
    uint outW, uint outH,
    global char* output)
{
    const float visibleThreshold = 15;

    uint x = get_global_id(0);
    uint y = get_global_id(1);

    if (x >= outW || y >= outH) {
        return;
    }

    float2 normCoord = (float2)(
        (float)x / (float)outW,
        (float)y / (float)outH
    );

    float4 pixel = read_imagef(inputImg, smp, normCoord);
    float lum = (299*pixel.r + 587*pixel.g + 114*pixel.b) / 1000;
    char is_visible = lum > visibleThreshold;

    output[y*outW+x] = is_visible;
}
