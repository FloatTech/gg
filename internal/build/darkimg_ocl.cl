kernel void isdark(
    read_only image2d_t inputImg,
    sampler_t smp,
    uint outW, uint outH,
    global int* visibleCount)
{
    const uint visibleThreshold = 15;

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
    uint lum = (299*pixel.x + 587*pixel.y + 114*pixel.z) *256 / 1000;

    int is_visible = (lum > visibleThreshold)?1:0;
    if (is_visible) {
        atomic_add((volatile __global int*)visibleCount, is_visible);
    }
}
