kernel void isdark(
    read_only image2d_t inputImg,
    sampler_t smp,
    uint outW, uint outH,
    global int* visibleCount)
{
    const uint visibleThreshold = 15;

    uint x = get_global_id(0);
    uint y = get_global_id(1);

    uint lid = get_local_id(1) * get_local_size(0) + get_local_id(0);

    local int localVisible;
    if (lid == 0) {
        localVisible = 0;
    }
    barrier(CLK_LOCAL_MEM_FENCE);

    if (x < outW && y < outH) {
        float2 normCoord = (float2)(
            (float)x / (float)outW,
            (float)y / (float)outH
        );

        float4 pixel = read_imagef(inputImg, smp, normCoord);
        float flum = mad(0.299f, pixel.x, mad(0.587f, pixel.y, 0.114f * pixel.z));

        if (flum > (15.0f / 256.0f)) {
            atomic_add(&localVisible, 1);
        }
    }

    barrier(CLK_LOCAL_MEM_FENCE);
    if (lid == 0 && localVisible > 0) {
        atomic_add(visibleCount, localVisible);
    }
}
