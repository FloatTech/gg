; ModuleID = 'device_bezier_kern.bc'
source_filename = "bezier_sycl.cpp"
target datalayout = "e-i64:64-v16:16-v24:32-v32:32-v48:64-v96:128-v192:256-v256:256-v512:512-v1024:1024-n8:16:32:64-G1"
target triple = "spirv64-unknown-unknown"

%struct.point = type { double, double }

@__spirv_BuiltInGlobalInvocationId = external local_unnamed_addr addrspace(1) constant <3 x i64>, align 32

; Function Attrs: mustprogress nofree norecurse nosync nounwind willreturn memory(argmem: write, inaccessiblemem: write)
define spir_kernel void @__sycl_kernel_quadratic(double noundef %0, double noundef %1, double noundef %2, double noundef %3, double noundef %4, double noundef %5, double noundef %6, ptr addrspace(1) noundef writeonly align 8 captures(none) %7) local_unnamed_addr #0 !kernel_arg_buffer_location !8 !sycl_used_aspects !9 !sycl_fixed_targets !11 !sycl_kernel_omit_args !12 {
  %9 = load i64, ptr addrspace(1) @__spirv_BuiltInGlobalInvocationId, align 32, !noalias !13
  %10 = icmp ult i64 %9, 2147483648
  tail call void @llvm.assume(i1 %10)
  %11 = uitofp nneg i64 %9 to double
  %12 = fdiv double %11, %6
  %13 = fsub double 1.000000e+00, %12
  %14 = fmul double %13, %13
  %15 = fmul double %13, 2.000000e+00
  %16 = fmul double %15, %12
  %17 = fmul double %12, %12
  %18 = fmul double %16, %2
  %19 = tail call double @llvm.fmuladd.f64(double %14, double %0, double %18)
  %20 = tail call double @llvm.fmuladd.f64(double %17, double %4, double %19)
  %21 = getelementptr inbounds %struct.point, ptr addrspace(1) %7, i64 %9
  store double %20, ptr addrspace(1) %21, align 8
  %22 = fmul double %16, %3
  %23 = tail call double @llvm.fmuladd.f64(double %14, double %1, double %22)
  %24 = tail call double @llvm.fmuladd.f64(double %17, double %5, double %23)
  %25 = getelementptr inbounds i8, ptr addrspace(1) %21, i64 8
  store double %24, ptr addrspace(1) %25, align 8
  ret void
}

; Function Attrs: nocallback nofree nosync nounwind speculatable willreturn memory(none)
declare !sycl_used_aspects !9 double @llvm.fmuladd.f64(double, double, double) #1

; Function Attrs: nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: write)
declare void @llvm.assume(i1 noundef) #2

; Function Attrs: mustprogress norecurse nounwind
define spir_kernel void @__sycl_kernel_cubic(double noundef %0, double noundef %1, double noundef %2, double noundef %3, double noundef %4, double noundef %5, double noundef %6, double noundef %7, double noundef %8, ptr addrspace(1) noundef align 8 %9) local_unnamed_addr #3 !kernel_arg_buffer_location !20 !sycl_used_aspects !9 !sycl_fixed_targets !11 !sycl_kernel_omit_args !21 {
  %11 = addrspacecast ptr addrspace(1) %9 to ptr addrspace(4)
  tail call spir_func void @cubic(double noundef %0, double noundef %1, double noundef %2, double noundef %3, double noundef %4, double noundef %5, double noundef %6, double noundef %7, double noundef %8, ptr addrspace(4) noundef %11) #5
  ret void
}

; Function Attrs: mustprogress norecurse nounwind
define linkonce_odr spir_func void @cubic(double noundef %0, double noundef %1, double noundef %2, double noundef %3, double noundef %4, double noundef %5, double noundef %6, double noundef %7, double noundef %8, ptr addrspace(4) noundef %9) local_unnamed_addr #4 !sycl_used_aspects !9 {
  %11 = load i64, ptr addrspace(1) @__spirv_BuiltInGlobalInvocationId, align 32, !noalias !22
  %12 = icmp ult i64 %11, 2147483648
  tail call void @llvm.assume(i1 %12)
  %13 = uitofp nneg i64 %11 to double
  %14 = fdiv double %13, %8
  %15 = fsub double 1.000000e+00, %14
  %16 = fmul double %15, %15
  %17 = fmul double %16, %15
  %18 = fmul double %15, 3.000000e+00
  %19 = fmul double %18, %15
  %20 = fmul double %19, %14
  %21 = fmul double %18, %14
  %22 = fmul double %21, %14
  %23 = fmul double %14, %14
  %24 = fmul double %23, %14
  %25 = fmul double %20, %2
  %26 = tail call double @llvm.fmuladd.f64(double %17, double %0, double %25)
  %27 = tail call double @llvm.fmuladd.f64(double %22, double %4, double %26)
  %28 = tail call double @llvm.fmuladd.f64(double %24, double %6, double %27)
  %29 = getelementptr inbounds nuw %struct.point, ptr addrspace(4) %9, i64 %11
  store double %28, ptr addrspace(4) %29, align 8
  %30 = fmul double %20, %3
  %31 = tail call double @llvm.fmuladd.f64(double %17, double %1, double %30)
  %32 = tail call double @llvm.fmuladd.f64(double %22, double %5, double %31)
  %33 = tail call double @llvm.fmuladd.f64(double %24, double %7, double %32)
  %34 = getelementptr inbounds nuw i8, ptr addrspace(4) %29, i64 8
  store double %33, ptr addrspace(4) %34, align 8
  ret void
}

declare dso_local spir_func i32 @_Z18__spirv_ocl_printfPU3AS2Kcz(ptr addrspace(2), ...)

attributes #0 = { mustprogress nofree norecurse nosync nounwind willreturn memory(argmem: write, inaccessiblemem: write) "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "sycl-module-id"="bezier_sycl.cpp" "sycl-nd-range-kernel"="1" "sycl-optlevel"="2" "uniform-work-group-size"="true" }
attributes #1 = { nocallback nofree nosync nounwind speculatable willreturn memory(none) }
attributes #2 = { nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: write) }
attributes #3 = { mustprogress norecurse nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "sycl-module-id"="bezier_sycl.cpp" "sycl-nd-range-kernel"="1" "sycl-optlevel"="2" "uniform-work-group-size"="true" }
attributes #4 = { mustprogress norecurse nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "sycl-nd-range-kernel"="1" "sycl-optlevel"="2" }
attributes #5 = { nounwind }

!llvm.linker.options = !{!0, !1}
!llvm.module.flags = !{!2, !3, !4}
!opencl.spir.version = !{!5}
!spirv.Source = !{!6}
!llvm.ident = !{!7}

!0 = !{!"-llibcpmt"}
!1 = !{!"/alternatename:_Avx2WmemEnabled=_Avx2WmemEnabledWeakValue"}
!2 = !{i32 1, !"wchar_size", i32 2}
!3 = !{i32 1, !"sycl-device", i32 1}
!4 = !{i32 7, !"frame-pointer", i32 2}
!5 = !{i32 1, i32 2}
!6 = !{i32 4, i32 100000}
!7 = !{!"clang version 21.0.0git (https://github.com/intel/llvm d5f649b706f63b5c74e1929bc95db8de91085560)"}
!8 = !{i32 -1, i32 -1, i32 -1, i32 -1, i32 -1, i32 -1, i32 -1, i32 -1}
!9 = !{!10}
!10 = !{!"fp64", i32 6}
!11 = !{}
!12 = !{i1 false, i1 false, i1 false, i1 false, i1 false, i1 false, i1 false, i1 false}
!13 = !{!14, !16, !18}
!14 = distinct !{!14, !15, !"_ZN7__spirv29InitSizesSTGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEE8initSizeEv: argument 0"}
!15 = distinct !{!15, !"_ZN7__spirv29InitSizesSTGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEE8initSizeEv"}
!16 = distinct !{!16, !17, !"_ZN7__spirv22initGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEEET0_v: argument 0"}
!17 = distinct !{!17, !"_ZN7__spirv22initGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEEET0_v"}
!18 = distinct !{!18, !19, !"_ZNK4sycl3_V17nd_itemILi1EE13get_global_idEv: argument 0"}
!19 = distinct !{!19, !"_ZNK4sycl3_V17nd_itemILi1EE13get_global_idEv"}
!20 = !{i32 -1, i32 -1, i32 -1, i32 -1, i32 -1, i32 -1, i32 -1, i32 -1, i32 -1, i32 -1}
!21 = !{i1 false, i1 false, i1 false, i1 false, i1 false, i1 false, i1 false, i1 false, i1 false, i1 false}
!22 = !{!23, !25, !27}
!23 = distinct !{!23, !24, !"_ZN7__spirv29InitSizesSTGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEE8initSizeEv: argument 0"}
!24 = distinct !{!24, !"_ZN7__spirv29InitSizesSTGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEE8initSizeEv"}
!25 = distinct !{!25, !26, !"_ZN7__spirv22initGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEEET0_v: argument 0"}
!26 = distinct !{!26, !"_ZN7__spirv22initGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEEET0_v"}
!27 = distinct !{!27, !28, !"_ZNK4sycl3_V17nd_itemILi1EE13get_global_idEv: argument 0"}
!28 = distinct !{!28, !"_ZNK4sycl3_V17nd_itemILi1EE13get_global_idEv"}
